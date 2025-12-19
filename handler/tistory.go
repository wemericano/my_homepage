package handler

import (
	"context"
	"fmt"
	"my-homepage/config"
	"my-homepage/tistory"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadToTistory 티스토리에 글 업로드
func UploadToTistory(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FAIL",
			"message": "요청 데이터가 올바르지 않습니다.",
		})
		return
	}

	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FAIL",
			"message": "제목과 내용은 필수입니다.",
		})
		return
	}

	// 설정 로드
	cfg := config.LoadConfig()

	if cfg.Tistory.Email == "" || cfg.Tistory.Password == "" || cfg.Tistory.BlogName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": "티스토리 설정이 올바르지 않습니다. 환경 변수를 확인해주세요.",
		})
		return
	}

	// 티스토리 클라이언트 생성
	client := tistory.NewClient(
		cfg.Tistory.Email,
		cfg.Tistory.Password,
		cfg.Tistory.BlogName,
		cfg.Tistory.Headless,
	)
	defer client.Close()

	// 컨텍스트 생성 (타임아웃 설정)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 글 작성
	fmt.Printf("📝 티스토리 글 작성 시작: %s\n", req.Title)
	result, err := client.WritePost(ctx, req.Title, req.Content)
	if err != nil {
		fmt.Printf("❌ 티스토리 업로드 실패: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": fmt.Sprintf("티스토리 업로드 실패: %v", err),
		})
		return
	}

	fmt.Printf("✅ 티스토리 업로드 성공: %s\n", result.URL)

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "티스토리 업로드가 완료되었습니다.",
		"data": gin.H{
			"postId": result.PostID,
			"url":    result.URL,
		},
	})
}
