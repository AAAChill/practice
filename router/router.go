package router

import (
	"github.com/gin-gonic/gin"
	"practice/api"
	"practice/middleware"
)

func InitRouter() *gin.Engine {
	g := gin.Default()
	g.GET("/metrics", api.Metrics)

	v1 := g.Group("/v1")
	v1.GET("/to_original/:short_link", api.Redirect)
	// 使用桶限流
	v1.Use(middleware.TokenBucketRedis())
	v1.POST("/shorten", api.Shorten)
	return g
}
