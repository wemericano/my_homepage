package main

import (
	"log"

	"my-homepage/config"
	db "my-homepage/database"
	"my-homepage/router"
	"my-homepage/scheduler"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
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

	// 자동 블로그 생성 스케줄러 시작
	batchScheduler := scheduler.NewBatchScheduler()
	batchScheduler.Start()
	defer batchScheduler.Stop()
	log.Println("[MAIN] SCHEDULER STARTED")

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
	log.Println("[MAIN] DB 연결 시도 중...")
	log.Printf("[MAIN] DB 설정 - Host: %s, Port: %s, User: %s, DBName: %s",
		dbConfig.Address, dbConfig.Port, dbConfig.User, dbConfig.Scheme)

	// DB 설정이 비어있으면 연결 시도 건너뛰기
	if dbConfig.Address == "" || dbConfig.Port == "" {
		log.Println("[MAIN] DB 설정이 비어있어 연결을 건너뜁니다.")
	} else {
		database, err := db.Open(dbConfig)
		if err != nil {
			log.Printf("[MAIN] DB 연결 실패 (계속 진행): %v", err)
		} else {
			defer database.Close()
			log.Println("[MAIN] DB SUCCESS")
		}
	}

	// 서버 실행
	if err := r.Run(":" + config.Server.Port); err != nil {
		log.Fatal(err)
	}
	log.Println("[MAIN] SERVER ON SUCCESS")
}
