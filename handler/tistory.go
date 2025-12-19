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
}
