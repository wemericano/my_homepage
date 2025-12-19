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
		BlogType     string `json:"blogType"`
		MainCategory string `json:"mainCategory"`
		SubCategory  string `json:"subCategory"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "잘못된 요청입니다."})
		return
	}

	// blogType 검증
	if request.BlogType != "hot" && request.BlogType != "news" && request.BlogType != "sports" && request.BlogType != "lotto" && request.BlogType != "fortune" && request.BlogType != "custom" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "잘못된 블로그 타입입니다."})
		return
	}

	// 직접입력일 때만 카테고리 검증
	if request.BlogType == "custom" {
		if request.MainCategory == "" || request.SubCategory == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "큰제목과 작은제목을 모두 선택해주세요."})
			return
		}
	}

	// 설정 불러오기
	cfg := config.LoadConfig()
	if cfg.GPT.APIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": "GPT API 키가 설정되지 않았습니다.",
			"data": gin.H{
				"title":   getDefaultTitle(request.BlogType, request.MainCategory, request.SubCategory),
				"content": "GPT API 키를 .env 파일에 GPT_API_KEY로 설정해주세요.",
			},
		})
		return
	}

	// GPT API 호출
	title, content, err := callGPTAPI(cfg.GPT.APIKey, request.BlogType, request.MainCategory, request.SubCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": "GPT API 호출 중 오류가 발생했습니다: " + err.Error(),
			"data": gin.H{
				"title":   getDefaultTitle(request.BlogType, request.MainCategory, request.SubCategory),
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

// 기본 제목 생성 함수
func getDefaultTitle(blogType, mainCategory, subCategory string) string {
	switch blogType {
	case "hot":
		return "Hot 블로그"
	case "news":
		return "News 블로그"
	case "sports":
		return "스포츠 블로그"
	case "lotto":
		return "로또 블로그"
	case "fortune":
		return "운세 블로그"
	case "custom":
		return mainCategory + " - " + subCategory
	default:
		return "블로그"
	}
}

// GPT API 호출 함수
func callGPTAPI(apiKey, blogType, mainCategory, subCategory string) (string, string, error) {
	var prompt string

	switch blogType {
	case "hot":
		prompt = `
					현재 가장 핫한 내용을 작성해
					요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
					내용은 500자 이상으로 작성해주세요.
					적절한 이모티콘/이모지를 사용하여 가독성 향상
					섹션 구분을 위해 이모티콘 활용
					독자들이 읽기 쉽고 흥미롭게 작성
				`
	case "news":
		prompt = `
					지금 웹에서 뉴스들 검색해보고 가장 핫한 뉴스에 관해 블로그를 만들어
					요새 핫한 블로그 처럼 제목과 내용을 작성해
					내용은 500자 이상으로 작성해주세요.
					적절한 이모티콘/이모지를 사용하여 가독성 향상
					섹션 구분을 위해 이모티콘 활용
					독자들이 읽기 쉽고 흥미롭게 작성
				`
	case "sports":
		prompt = `
					오늘 NBA, KBO, MLB 등 주요 스포츠 리그의 최신 소식과 뉴스를 검색해서 블로그를 작성해주세요.
					경기 결과, 선수 소식, 트레이드 뉴스, 주요 이슈 등을 다루어주세요.
					요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
					내용은 500자 이상으로 작성해주세요.
					적절한 이모티콘/이모지를 사용하여 가독성 향상 (예: ⚽, 🏀, ⚾, 🏈, 🎯 등)
					섹션 구분을 위해 이모티콘 활용
					독자들이 읽기 쉽고 흥미롭게 작성
				`
	case "lotto":
		prompt = `
					로또 당첨번호 분석, 번호 추천, 통계, 꿈 해몽, 행운의 숫자 등 로또 관련 내용으로 블로그를 작성해주세요.
					최근 당첨번호 패턴 분석이나 행운의 번호 추천 등을 포함해주세요.
					요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
					내용은 500자 이상으로 작성해주세요.
					적절한 이모티콘/이모지를 사용하여 가독성 향상 (예: 🎰, 🍀, ⭐, 💰, 🎯 등)
					섹션 구분을 위해 이모티콘 활용
					독자들이 읽기 쉽고 흥미롭게 작성
				`
	case "fortune":
		prompt = `
					오늘의 운세, 별자리 운세, 타로, 사주, 꿈 해몽 등 운세 관련 내용으로 블로그를 작성해주세요.
					오늘의 행운의 색상, 숫자, 방향, 조언 등을 포함해주세요.
					요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
					내용은 500자 이상으로 작성해주세요.
					적절한 이모티콘/이모지를 사용하여 가독성 향상 (예: 🔮, ✨, ⭐, 🌟, 🍀, 🔯 등)
					섹션 구분을 위해 이모티콘 활용
					독자들이 읽기 쉽고 흥미롭게 작성
				`
	case "custom":
		prompt = fmt.Sprintf(`
								다음 주제에 대한 블로그 포스트를 작성해주세요.

								카테고리: %s
								
								주제: %s

								제목과 내용을 포함하여 작성해주세요.
								내용은 500자 이상으로 작성해주세요.
								요새 핫한 블로그처럼 제목과 내용을 작성해주세요.
								적절한 이모티콘/이모지를 사용하여 가독성 향상
								섹션 구분을 위해 이모티콘 활용
								독자들이 읽기 쉽고 흥미롭게 작성
							`, mainCategory, subCategory)
	default:
		return "", "", fmt.Errorf("지원하지 않는 블로그 타입입니다: %s", blogType)
	}

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

	// 제목과 내용 분리
	lines := strings.Split(responseText, "\n")
	var title string
	var content string

	// 제목 추출 로직
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])

		// "제목:" 또는 "Title:" 형식인 경우
		if strings.HasPrefix(firstLine, "제목:") || strings.HasPrefix(firstLine, "Title:") {
			title = strings.TrimSpace(strings.SplitN(firstLine, ":", 2)[1])
			content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		} else {
			// 첫 줄이 제목일 가능성이 높음 (짧고 명확한 경우)
			// 첫 줄이 너무 길면(50자 이상) 제목이 아닐 수 있음
			if len(firstLine) > 0 && len(firstLine) < 100 && !strings.Contains(firstLine, "```") {
				title = firstLine
				// 제목에 이모티콘이나 특수문자가 많으면 제목일 가능성 높음
				if len(lines) > 1 {
					content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
				} else {
					content = responseText
				}
			} else {
				// 첫 줄이 제목 같지 않으면 전체를 내용으로, GPT가 생성한 제목 찾기
				// "## " 또는 "# " 형식의 마크다운 제목 찾기
				foundTitle := false
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "# ") {
						title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
						foundTitle = true
						break
					} else if strings.HasPrefix(trimmed, "## ") {
						title = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
						foundTitle = true
						break
					}
				}

				if !foundTitle {
					// 제목을 찾지 못했으면 첫 줄을 제목으로 사용
					title = firstLine
					if len(lines) > 1 {
						content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
					} else {
						content = responseText
					}
				} else {
					// 제목을 찾았으면 나머지를 내용으로
					// 제목 라인 제거
					contentLines := []string{}
					titleRemoved := false
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if !titleRemoved && (strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ")) {
							titleRemoved = true
							continue
						}
						contentLines = append(contentLines, line)
					}
					content = strings.TrimSpace(strings.Join(contentLines, "\n"))
				}
			}
		}
	} else {
		// 빈 응답인 경우
		title = getDefaultTitle(blogType, mainCategory, subCategory)
		content = responseText
	}

	// 제목이 비어있으면 기본 제목 사용
	if strings.TrimSpace(title) == "" {
		title = getDefaultTitle(blogType, mainCategory, subCategory)
	}

	// 내용이 비어있으면 전체 응답 사용
	if strings.TrimSpace(content) == "" {
		content = responseText
	}

	return title, content, nil
}
