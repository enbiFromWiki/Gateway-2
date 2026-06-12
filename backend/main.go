package main

import (
	// "gateway/backend/app"
	"gateway/backend/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	// app.Run()
	r := gin.Default()

	r.GET("/login", auth.Login)
	r.GET("/auth/callback", auth.Login)
	r.GET("/call", auth.ApiTest)

	r.Run("127.0.0.1:8080")
}
