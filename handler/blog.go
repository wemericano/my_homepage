package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-homepage/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GPT API 요청 구조체
type GPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

// 블로그 생성 API
func GenerateBlog(c *gin.Context) {
	var request struct {
		MainCategory string `json:"mainCategory"`
		SubCategory  string `json:"subCategory"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "잘못된 요청입니다."})
		return
	}

	// 설정 불러오기
	cfg := config.LoadConfig()
	if cfg.GPT.APIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": "GPT API 키가 설정되지 않았습니다.",
			"data": gin.H{
				"title":   request.MainCategory + " - " + request.SubCategory,
				"content": "GPT API 키를 .env 파일에 GPT_API_KEY로 설정해주세요.",
			},
		})
		return
	}

	// GPT API 호출
	title, content, err := callGPTAPI(cfg.GPT.APIKey, request.MainCategory, request.SubCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": "GPT API 호출 중 오류가 발생했습니다: " + err.Error(),
			"data": gin.H{
				"title":   request.MainCategory + " - " + request.SubCategory,
				"content": "GPT API 호출 실패: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "블로그 생성이 완료되었습니다.",
		"data": gin.H{
			"title":   title,
			"content": content,
		},
	})
}

// GPT API 호출 함수
func callGPTAPI(apiKey, mainCategory, subCategory string) (string, string, error) {
	prompt := fmt.Sprintf(`
							다음 주제에 대한 블로그 포스트를 작성해주세요.

							카테고리: %s
							주제: %s

							제목과 내용을 포함하여 작성해주세요.
							내용은 500자 이상으로 작성해주세요.
							요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
						`, mainCategory, subCategory)

	gptReq := GPTRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(gptReq)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GPT API 오류: %s", string(body))
	}

	var gptResp GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return "", "", err
	}

	if len(gptResp.Choices) == 0 {
		return "", "", fmt.Errorf("GPT 응답에 내용이 없습니다")
	}

	// 응답에서 제목과 내용 추출
	responseText := gptResp.Choices[0].Message.Content

	// 제목과 내용 분리 (첫 줄을 제목으로, 나머지를 내용으로)
	lines := strings.Split(responseText, "\n")
	title := mainCategory + " - " + subCategory
	content := responseText

	// 제목이 명확히 구분되어 있으면 추출
	if len(lines) > 0 && strings.HasPrefix(lines[0], "제목:") || strings.HasPrefix(lines[0], "Title:") {
		title = strings.TrimSpace(strings.Split(lines[0], ":")[1])
		content = strings.Join(lines[1:], "\n")
	}

	return title, content, nil
}
