package router

import (
	"my-homepage/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// public 폴더 정적 파일 제공
	r.Static("/public", "./public")
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	api := r.Group("/api")

	// 회원가입 API
	api.POST("/signup", service.AddSignup)
}
