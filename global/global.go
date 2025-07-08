package global

import (
	"github.com/redis/go-redis/v9"
	"math/rand"
	"sync"
	"time"
)

var (
	TokenBucket chan int
	RedisClient *redis.Client
	IPBucketMap *sync.Map
)

func Init() {
	rand.Seed(time.Now().UnixNano())
	go TokenBucketFiller(5)
	go IPTokenBucketFiller(5)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// TokenBucketFiller
// @Description 令牌桶填充
func TokenBucketFiller(maxToken int) {
	//初始化一个令牌桶来存储令牌
	TokenBucket = make(chan int, maxToken)
	tick := time.Tick(time.Second / time.Duration(maxToken))

	for {
		select {
		case TokenBucket <- 1:
		default:
		}
		<-tick
	}
}

func IPTokenBucketFiller(maxToken int) {
	IPBucketMap = new(sync.Map)
	tick := time.Tick(time.Second)
	for {
		IPBucketMap.Range(func(key, value any) bool {
			select {
			case value.(chan int) <- 1:
			default:
			}
			return true
		})

		<-tick
	}
}
