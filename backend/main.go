package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mvlberry/backend/api/handlers"
)

func setCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(setCORS())

	//if /api/signup is hit with a POST request, run the signup function
	r.POST("/api/signup", handlers.Signup)

	//if /api/login is hit with a POST request, run the login function
	r.POST("/api/login", handlers.Login)

	r.GET("api/home/getfriends", handlers.Getfriends)

	r.GET("api/ws", handlers.Startws)

	r.Run("localhost:8080")
}
