package router

import (
	"my-homepage/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// public 폴더 정적 파일 제공
	r.Static("/public", "./public")
	r.Static("/page", "./public/page")
	r.Static("/js", "./public/js")
	r.Static("/assets", "./public/assets")
	r.Static("/base", "./public/base")

	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	api := r.Group("/api")

	// 회원가입
	api.POST("/signup", service.AddSignup)

	// 로그인
	api.POST("/login", service.Login)

	// lotto
	api.POST("/lotto", service.GetLottoList)
}
