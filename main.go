package main

import (
	"practice/global"
	"practice/router"
)

func main() {
	global.Init()
	g := router.InitRouter()
	g.Run(":8080")
}
