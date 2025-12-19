package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 카카오 로그인 버튼 클릭까지만
func UploadToTistory(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "티스토리 업로드 요청이 처리되었습니다.",
	})
}
