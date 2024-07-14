package main

import (
	"SkillSwap/internal"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Allow-Origin", origin)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.Use(CORSMiddleware())
	internal.SetRouters(router)
	go internal.SendPing()
	go internal.SaveMessagesToDB()
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
	}

}
