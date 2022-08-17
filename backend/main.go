package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mvlberry/backend/api/handlers"
)

func main() {
	r := gin.Default()

	//if /api/signup is hit with a POST request, run the signup function
	r.POST("/api/signup", handlers.Signup)

	r.Run("localhost:8080")
}
