package main

import (
	"douyin/repository"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	go RunMessageServer()

	if err := repository.InitDB(); err != nil {
		os.Exit(-1)
	}
	if err := repository.InitRedis(); err != nil {
		os.Exit(-1)
	}
	if err := repository.InitKafka(); err != nil {
		os.Exit(-1)
	}
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
