package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"practice/global"
	"strconv"
	"strings"
)

func Metrics(c *gin.Context) {
	tokenCountPre := "tokenCount-IP:"
	tokenKeys := global.RedisClient.Keys(context.Background(), tokenCountPre+"*").Val()
	returnStr := ""
	name := "token_bucket_ip_num"
	tag := "ip"
	for _, key := range tokenKeys {
		tokenCount, _ := global.RedisClient.Get(context.Background(), key).Int()
		ip := strings.TrimPrefix(key, tokenCountPre)
		returnStr += name + "{" + tag + "=\"" + ip + "\"}" + " " + strconv.Itoa(tokenCount) + "\n"
	}

	c.Set("Content-Type", "text/plain")
	c.String(200, returnStr)
}
