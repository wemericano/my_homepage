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

	// Chrome DevTools 요청 처리 (404 방지)
	r.GET("/.well-known/appspecific/com.chrome.devtools.json", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	api := r.Group("/api")

	// 회원가입
	api.POST("/signup", handler.AddSignup)

	// 로그인
	api.POST("/login", handler.Login)

	// lotto
	api.POST("/lotto", handler.GetLottoList)
	api.POST("/analyze/v1", handler.AnalyzeV1)
	api.POST("/analyze/v2", handler.AnalyzeV2)

	// blog
	api.POST("/blog/generate", handler.GenerateBlog)
}
