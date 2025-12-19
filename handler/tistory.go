package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 티스토리 업로드 API
func UploadToTistory(c *gin.Context) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("[TISTORY] 요청 바인딩 실패: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "잘못된 요청입니다."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "티스토리 업로드 요청이 처리되었습니다.",
	})
}
