package handler

import (
	"context"
	"fmt"
	"my-homepage/config"
	"my-homepage/tistory"
	"net/http"
	"regexp"
	"strings"
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

	// 마크다운을 HTML로 변환
	htmlContent := ConvertMarkdownToHTML(req.Content)

	// 글 작성
	fmt.Printf("📝 티스토리 글 작성 시작: %s\n", req.Title)
	result, err := client.WritePost(ctx, req.Title, htmlContent)
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

// ConvertMarkdownToHTML 마크다운을 티스토리에 맞는 HTML로 변환 (외부에서 사용 가능하도록 export)
func ConvertMarkdownToHTML(markdown string) string {
	content := markdown

	// 헤더 변환 (# ## ###) - 먼저 처리
	content = regexp.MustCompile(`(?m)^###\s+(.+)$`).ReplaceAllStringFunc(content, func(match string) string {
		text := regexp.MustCompile(`^###\s+(.+)$`).FindStringSubmatch(match)[1]
		return fmt.Sprintf(`<h3 style="font-size: 1.3em; font-weight: bold; margin-top: 20px; margin-bottom: 10px; color: #333; border-left: 4px solid #4CAF50; padding-left: 10px;">%s</h3>`, htmlEscape(text))
	})
	content = regexp.MustCompile(`(?m)^##\s+(.+)$`).ReplaceAllStringFunc(content, func(match string) string {
		text := regexp.MustCompile(`^##\s+(.+)$`).FindStringSubmatch(match)[1]
		return fmt.Sprintf(`<h2 style="font-size: 1.5em; font-weight: bold; margin-top: 25px; margin-bottom: 15px; color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 8px;">%s</h2>`, htmlEscape(text))
	})
	content = regexp.MustCompile(`(?m)^#\s+(.+)$`).ReplaceAllStringFunc(content, func(match string) string {
		text := regexp.MustCompile(`^#\s+(.+)$`).FindStringSubmatch(match)[1]
		return fmt.Sprintf(`<h1 style="font-size: 1.8em; font-weight: bold; margin-top: 30px; margin-bottom: 20px; color: white; text-align: center; padding: 15px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 8px;">%s</h1>`, htmlEscape(text))
	})

	// 리스트 변환 (- 또는 *)
	lines := strings.Split(content, "\n")
	var result []string
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 이미 헤더 태그인 경우 그대로 유지
		if strings.HasPrefix(trimmed, "<h") {
			if inList {
				result = append(result, `</ul>`)
				inList = false
			}
			result = append(result, line)
			continue
		}

		// 리스트 항목인지 확인
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList {
				result = append(result, `<ul style="list-style-type: none; padding-left: 0; margin: 15px 0;">`)
				inList = true
			}
			item := strings.TrimPrefix(trimmed, "- ")
			item = strings.TrimPrefix(item, "* ")
			// 강조 표시 처리
			item = processInlineFormatting(item)
			result = append(result, fmt.Sprintf(`<li style="padding: 8px 0; padding-left: 25px; position: relative; line-height: 1.6;">✨ %s</li>`, item))
		} else {
			if inList {
				result = append(result, `</ul>`)
				inList = false
			}
			if trimmed != "" {
				// 일반 텍스트도 인라인 포맷팅 처리
				processed := processInlineFormatting(trimmed)
				result = append(result, fmt.Sprintf(`<p style="margin: 10px 0; line-height: 1.8;">%s</p>`, processed))
			} else {
				result = append(result, `<br>`)
			}
		}
	}

	if inList {
		result = append(result, `</ul>`)
	}

	content = strings.Join(result, "\n")

	// 문단 스타일링
	content = fmt.Sprintf(`<div style="font-family: 'Noto Sans KR', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.8; color: #333; max-width: 100%%; padding: 20px;">%s</div>`, content)

	return content
}

// htmlEscape HTML 특수 문자 이스케이프
func htmlEscape(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, `"`, "&quot;")
	return text
}

// processInlineFormatting 인라인 포맷팅 처리 (강조, 기울임)
func processInlineFormatting(text string) string {
	// 강조 표시 (**텍스트**) - 먼저 처리하고 임시로 치환
	strongPattern := regexp.MustCompile(`\*\*(.+?)\*\*`)
	strongPlaceholder := "___STRONG_PLACEHOLDER___"
	var strongMatches []string
	strongIndex := 0

	text = strongPattern.ReplaceAllStringFunc(text, func(match string) string {
		inner := strongPattern.FindStringSubmatch(match)[1]
		strongMatches = append(strongMatches, fmt.Sprintf(`<strong style="font-weight: bold; color: #e74c3c;">%s</strong>`, htmlEscape(inner)))
		placeholder := fmt.Sprintf("%s%d", strongPlaceholder, strongIndex)
		strongIndex++
		return placeholder
	})

	// 기울임 (*텍스트*) - **가 아닌 단일 *만 처리
	emPattern := regexp.MustCompile(`([^*]|^)\*([^*]+?)\*([^*]|$)`)
	emPlaceholder := "___EM_PLACEHOLDER___"
	var emMatches []string
	emIndex := 0

	text = emPattern.ReplaceAllStringFunc(text, func(match string) string {
		submatch := emPattern.FindStringSubmatch(match)
		before := submatch[1]
		inner := submatch[2]
		after := submatch[3]
		emMatches = append(emMatches, fmt.Sprintf(`<em style="font-style: italic; color: #9b59b6;">%s</em>`, htmlEscape(inner)))
		placeholder := fmt.Sprintf("%s%d", emPlaceholder, emIndex)
		emIndex++
		return before + placeholder + after
	})

	// 나머지 텍스트 이스케이프
	text = htmlEscape(text)

	// 강조 플레이스홀더 복원
	for i, match := range strongMatches {
		placeholder := fmt.Sprintf("%s%d", strongPlaceholder, i)
		text = strings.ReplaceAll(text, placeholder, match)
	}

	// 기울임 플레이스홀더 복원
	for i, match := range emMatches {
		placeholder := fmt.Sprintf("%s%d", emPlaceholder, i)
		text = strings.ReplaceAll(text, placeholder, match)
	}

	return text
}
