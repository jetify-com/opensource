package main

import (
	"go-gin/handlers/root"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.GET("/", root.Root)

	router.Run() // listen and serve on 0.0.0.0:8080
}
