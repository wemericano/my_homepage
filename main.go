package main

import (
	"log"

	"my-homepage/config"
	db "my-homepage/database"
	"my-homepage/router"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("[MAIN] SERVER STARTING")

	// 설정 불러오기
	config := config.LoadConfig()

	// 라우터 초기화
	// Gin 엔진 생성
	r := gin.Default()

	// 라우터 설정
	router.SetupRouter(r)
	log.Println("[MAIN] ROUTER SUCCESS")

	// DB 초기화
	dbConfig := db.DBConfig{
		Address:     config.DB.Host,
		Port:        config.DB.Port,
		User:        config.DB.User,
		Pw:          config.DB.Pass,
		Scheme:      config.DB.Name,
		MaxIdle:     10,
		MaxOpen:     100,
		MaxLifeTime: 10,
	}

	database, err := db.Open(dbConfig)
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	defer database.Close()

	log.Println("[MAIN] DB SUCCESS")

	// 서버 실행
	if err := r.Run(":" + config.Server.Port); err != nil {
		log.Fatal(err)
	}
	log.Println("[MAIN] SERVER ON SUCCESS")
}
