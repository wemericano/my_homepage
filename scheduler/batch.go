package scheduler

import (
	"context"
	"log"
	"my-homepage/config"
	"my-homepage/handler"
	"my-homepage/tistory"
	"time"

	"github.com/robfig/cron/v3"
)

// BatchScheduler 자동 블로그 생성 스케줄러
type BatchScheduler struct {
	cron *cron.Cron
}

// NewBatchScheduler 새로운 스케줄러 생성
func NewBatchScheduler() *BatchScheduler {
	// 타임존 설정 (한국 시간)
	kr, _ := time.LoadLocation("Asia/Seoul")
	c := cron.New(cron.WithLocation(kr), cron.WithSeconds())

	return &BatchScheduler{
		cron: c,
	}
}

// Start 스케줄러 시작
func (bs *BatchScheduler) Start() {
	// 오전 7시: 로또, 운세
	_, err := bs.cron.AddFunc("0 0 7 * * *", func() {
		bs.runBatch("lotto")
	})
	if err != nil {
		log.Printf("[SCHEDULER] 로또 스케줄 등록 실패: %v", err)
	} else {
		log.Println("[SCHEDULER] 로또 스케줄 등록 완료 (매일 오전 7시)")
	}

	_, err = bs.cron.AddFunc("0 0 7 * * *", func() {
		bs.runBatch("fortune")
	})
	if err != nil {
		log.Printf("[SCHEDULER] 운세 스케줄 등록 실패: %v", err)
	} else {
		log.Println("[SCHEDULER] 운세 스케줄 등록 완료 (매일 오전 7시)")
	}

	// 오후 13시: 스포츠
	_, err = bs.cron.AddFunc("0 0 13 * * *", func() {
		bs.runBatch("sports")
	})
	if err != nil {
		log.Printf("[SCHEDULER] 스포츠 스케줄 등록 실패: %v", err)
	} else {
		log.Println("[SCHEDULER] 스포츠 스케줄 등록 완료 (매일 오후 1시)")
	}

	// 4시간 간격으로 hot, news (0시, 4시, 8시, 12시, 16시, 20시)
	// hot: 0시, 8시, 16시
	_, err = bs.cron.AddFunc("0 0 0 * * *", func() {
		bs.runBatch("hot")
	})
	if err != nil {
		log.Printf("[SCHEDULER] Hot 스케줄 등록 실패 (0시): %v", err)
	} else {
		log.Println("[SCHEDULER] Hot 스케줄 등록 완료 (매일 0시)")
	}

	_, err = bs.cron.AddFunc("0 0 8 * * *", func() {
		bs.runBatch("hot")
	})
	if err != nil {
		log.Printf("[SCHEDULER] Hot 스케줄 등록 실패 (8시): %v", err)
	} else {
		log.Println("[SCHEDULER] Hot 스케줄 등록 완료 (매일 8시)")
	}

	_, err = bs.cron.AddFunc("0 0 16 * * *", func() {
		bs.runBatch("hot")
	})
	if err != nil {
		log.Printf("[SCHEDULER] Hot 스케줄 등록 실패 (16시): %v", err)
	} else {
		log.Println("[SCHEDULER] Hot 스케줄 등록 완료 (매일 16시)")
	}

	// news: 4시, 12시, 20시
	_, err = bs.cron.AddFunc("0 0 4 * * *", func() {
		bs.runBatch("news")
	})
	if err != nil {
		log.Printf("[SCHEDULER] News 스케줄 등록 실패 (4시): %v", err)
	} else {
		log.Println("[SCHEDULER] News 스케줄 등록 완료 (매일 4시)")
	}

	_, err = bs.cron.AddFunc("0 0 12 * * *", func() {
		bs.runBatch("news")
	})
	if err != nil {
		log.Printf("[SCHEDULER] News 스케줄 등록 실패 (12시): %v", err)
	} else {
		log.Println("[SCHEDULER] News 스케줄 등록 완료 (매일 12시)")
	}

	_, err = bs.cron.AddFunc("0 0 20 * * *", func() {
		bs.runBatch("news")
	})
	if err != nil {
		log.Printf("[SCHEDULER] News 스케줄 등록 실패 (20시): %v", err)
	} else {
		log.Println("[SCHEDULER] News 스케줄 등록 완료 (매일 20시)")
	}

	bs.cron.Start()
	log.Println("[SCHEDULER] 자동 블로그 생성 스케줄러 시작 완료")
}

// Stop 스케줄러 중지
func (bs *BatchScheduler) Stop() {
	bs.cron.Stop()
	log.Println("[SCHEDULER] 스케줄러 중지")
}

// runBatch 배치 작업 실행
func (bs *BatchScheduler) runBatch(blogType string) {
	log.Println("========================================")
	log.Printf("[BATCH] 자동 블로그 생성 시작 (타입: %s)", blogType)
	log.Printf("[BATCH] 날짜: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("========================================")

	cfg := config.LoadConfig()

	// GPT API 키 확인
	if cfg.GPT.APIKey == "" {
		log.Printf("[BATCH] ERROR: GPT API 키가 설정되지 않았습니다. (타입: %s)", blogType)
		return
	}

	// 티스토리 설정 확인
	if cfg.Tistory.Email == "" || cfg.Tistory.Password == "" || cfg.Tistory.BlogName == "" {
		log.Printf("[BATCH] ERROR: 티스토리 설정이 올바르지 않습니다. (타입: %s)", blogType)
		return
	}

	// 1. 블로그 생성 (handler/blog.go의 callGPTAPI 함수 사용)
	log.Printf("[BATCH] [1/2] 블로그 생성 중... (타입: %s)", blogType)
	title, content, err := handler.CallGPTAPI(cfg.GPT.APIKey, blogType, "", "")
	if err != nil {
		log.Printf("[BATCH] ERROR: 블로그 생성 실패 (타입: %s): %v", blogType, err)
		return
	}
	log.Printf("[BATCH] [SUCCESS] 블로그 생성 완료 - 제목: %s (타입: %s)", title, blogType)

	// 2. HTML 변환 (handler/tistory.go의 convertMarkdownToHTML 함수 사용)
	htmlContent := handler.ConvertMarkdownToHTML(content)

	// 3. 티스토리 업로드
	log.Println("[BATCH] [2/2] 티스토리 업로드 중...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := tistory.NewClient(
		cfg.Tistory.Email,
		cfg.Tistory.Password,
		cfg.Tistory.BlogName,
		cfg.Tistory.Headless,
	)
	defer client.Close()

	result, err := client.WritePost(ctx, title, htmlContent)
	if err != nil {
		log.Printf("[BATCH] ERROR: 티스토리 업로드 실패: %v", err)
		return
	}

	log.Printf("[BATCH] [SUCCESS] 티스토리 업로드 완료 - URL: %s (타입: %s)", result.URL, blogType)
	log.Println("========================================")
	log.Printf("[BATCH] 자동 블로그 생성 완료 (타입: %s)", blogType)
	log.Println("========================================")
}
