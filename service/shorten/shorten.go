package shorten

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"practice/global"
	"time"
)

func GetShortLink(url string) (short string, err error) {
	if url == "" {
		return "", errors.New("原始链接为空")
	}
	short, err = global.RedisClient.Get(context.Background(), "longURL:"+url).Result()
	if errors.Is(err, redis.Nil) {
		short = generateLink(6)
		global.RedisClient.Set(context.Background(), "longURL:"+url, short, time.Hour*24)
		global.RedisClient.Set(context.Background(), "shortURL:"+short, url, time.Hour*24)
	} else if err != nil {
		return "", err
	}
	return short, nil
}

// generateLink
// @Description 生成短链接
func generateLink(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		result := make([]byte, length)
		for i := range result {
			result[i] = charset[rand.Intn(len(charset))]
		}
		if isURLAvailable(string(result)) {
			return string(result)
		}
	}

}

// isURLAvailable
// @Description 检查短链接是否在缓存中
func isURLAvailable(short string) bool {
	_, err := global.RedisClient.Get(context.Background(), "shortURL:"+short).Result()
	if errors.Is(err, redis.Nil) {
		return true
	}
	return false
}

// GetLongURL
// @Description 获取原始链接
func GetLongURL(shortLink string) (string, error) {
	longURL, err := global.RedisClient.Get(context.Background(), "shortURL:"+shortLink).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("短链接不存在")
	} else if err != nil {
		return "", errors.New("获取失败")
	}
	return longURL, nil
}
