package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"practice/api"
	"practice/global"
	"sync"
	"time"
)

var limitQueue = map[string][]int64{}
var mutex = &sync.Mutex{}

// 限流中间件-滑动窗口-使用本地缓存
func RateLimitLocalCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		maxRate := 5

		ip := c.ClientIP()
		key := "rate_limit:" + ip

		mutex.Lock()
		defer mutex.Unlock()
		// 获取缓存的窗口队列
		queue, _ := limitQueue[key]

		// 去除时间范围外的缓存
		var newQueue []int64
		if len(queue) > 0 {
			for _, v := range queue {
				if time.Now().UnixMilli()-v < 1000 {
					newQueue = append(newQueue, v)
				}
			}
		}

		// 判断是否超出限制
		if len(newQueue) > maxRate-1 {
			api.ErrorResponse(c, 429, "too many requests", nil)
			c.Abort()
			return
		}

		// 更新缓存
		newQueue = append(newQueue, time.Now().UnixMilli())
		limitQueue[key] = newQueue
		c.Next()
	}
}

// 限流中间件-令牌桶-使用本地缓存
func TokenBucketLocalCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("桶长度:", len(global.TokenBucket), "桶容量:", cap(global.TokenBucket))
		select {
		case <-global.TokenBucket:
			c.Next()
		default:
			api.ErrorResponse(c, 429, "too many requests", nil)
			c.Abort()
			return
		}
	}
}

// 限流中间件-令牌桶-根据IP区分
func TokenBucketByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 本地多IP测试
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.ClientIP()
		}
		bucket, ok := global.IPBucketMap.Load(ip)
		if !ok {
			ipBucket := make(chan int, 5)
			for i := 0; i < 5; i++ {
				ipBucket <- 1
			}
			global.IPBucketMap.Store(ip, ipBucket)
			bucket = ipBucket
		}
		fmt.Println("IP:", ip, "令牌数量:", len(bucket.(chan int)))
		select {
		case <-bucket.(chan int):
			c.Next()
		default:
			api.ErrorResponse(c, 429, "too many requests", nil)
			c.Abort()
			return
		}
	}
}

// 限流中间件-滑动窗口-使用Redis
func RateLimitRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.FullPath()

		key := fmt.Sprintf("RateLimit:%s/%s", ip, path)
		res, err := global.RedisClient.Eval(
			context.Background(), RateLimit,
			[]string{key},
			time.Now().UnixMilli(), 1000, 5,
		).Int64()
		if err != nil {
			// 处理错误，先放空，默认拦截
			c.Abort()
			return
		}
		if res == 0 {
			api.ErrorResponse(c, 429, "too many requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 限流中间件-IP令牌桶-使用Redis
func TokenBucketRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		tokenKey := "tokenCount" + "-IP:" + ip
		lastTimeKey := "lastTime" + "-IP:" + ip
		res, err := global.RedisClient.Eval(
			context.Background(), TokenBucket,
			[]string{tokenKey, lastTimeKey},
			time.Now().UnixMilli(), 10, 1000,
		).Int64()
		if err != nil {
			// 处理错误，先放空，默认拦截
			log.Printf("[RateLimit] Redis Eval error: %v", err)
			c.Abort()
			return
		}
		if res == 0 {
			api.ErrorResponse(c, 429, "too many requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
