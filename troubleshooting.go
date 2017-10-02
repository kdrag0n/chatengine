package main

import (
	"chatengine/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func hubEntry(dest, name, desc string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s" title="%s">%s</a> - %s`, dest, name, name, desc))
}

func setupTroubleshooting(bot *ChatBot, router *gin.Engine) {
	group := router.Group("/troubleshooting")
	{
		group.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "hub.tmpl", gin.H{
				"title":   "Troubleshooting",
				"abspath": "troubleshooting/",
				"desc": template.HTML(`If you ever have problems with any parts of this website, this is the place to go.
You can try various tests to figure out and resolve your issue.
If you are unable to diagnose the issue after trying everything on this page, please contact the creator.
<br><br>
You will find links to various tests below.`),
				"listhead": "Tests:",
				"list": []template.HTML{
					hubEntry("time", "Time Check", "Check the accuracy of your time."),
					hubEntry("latency", "Latency Check", "Check the time taken for you to receive the server's response to a request."),
				},
			})
		})

		group.GET("/time", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "tb_time.tmpl", gin.H{
				"page_served": util.CurrentTimeMillis(),
			})
		})

		group.GET("/ctime_ms", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "%d", util.CurrentTimeMillis())
		})

		group.GET("/latency", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "tb_latency.tmpl", gin.H{
				"page_served": util.CurrentTimeMillis(),
			})
		})

		group.GET("/ping", func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		})
	}
}
