package main

import (
	"net/http"

	"chatengine/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

func setupManagement(bot *ChatBot, router *gin.Engine) {
	group := router.Group("/manage")
	{
		group.GET("/dashboard", func(ctx *gin.Context) {
			session := sessions.Get(ctx)
			if v := session.Get("logged_in"); v == nil || !v.(bool) {
				ctx.Redirect(http.StatusTemporaryRedirect, "/login")
				return
			}

			ctx.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
				"uname":      "Someone",
				"uid":        session.Get("user_id").(string),
				"ltime":      session.Get("logged_in_at").(time.Time).Format(generalTimeFmt),
				"loginIP":    session.Get("login_ip").(string),
				"rememberMe": session.Get("remember_me").(bool),
			})
		})
	}

	router.GET("/dashboard", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/manage/dashboard")
	})
}
