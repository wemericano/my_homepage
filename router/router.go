package router

import (
	"my-homepage/handler"

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
	api.POST("/signup", handler.AddSignup)

	// 로그인
	api.POST("/login", handler.Login)

	// lotto
	api.POST("/lotto", handler.GetLottoList)
	api.POST("/analyze/v1", handler.AnalyzeV1)
}
