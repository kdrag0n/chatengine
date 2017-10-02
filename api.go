package main

import (
	"github.com/gin-gonic/gin"
)

func setupAPI(bot *ChatBot, router *gin.Engine) {
	group := router.Group("/api", func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "deny")
		ctx.Next()
	})
	{
		group.POST("/ask", bot.HandleHTTP)
	}
}
