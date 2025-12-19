package config

import (
	"os"

	"github.com/joho/godotenv"
)

// 전체 Config 구조체
type Config struct {
	Server  ServerConfig
	DB      DBConfig
	GPT     GPTConfig
	Tistory TistoryConfig
}

// 서버 관련 설정
type ServerConfig struct {
	Port string
}

// DB 관련 설정
type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

// GPT 관련 설정
type GPTConfig struct {
	APIKey string
}

// Tistory 관련 설정
type TistoryConfig struct {
	Email    string
	Password string
	BlogName string
	Headless bool
}

// 메인 로드 함수
func LoadConfig() *Config {
	env := os.Getenv("ENV")
	isRender := os.Getenv("RENDER") != ""

	// 로컬에서 ENV=local 이거나, RENDER 환경변수가 없을 경우 .env 로드
	if env == "local" || !isRender {
		err := godotenv.Load()
		if err != nil {
			println("로컬 환경: .env 파일을 불러올 수 없습니다")
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SV_PORT", "0126"),
		},
		DB: DBConfig{
			User: getEnv("DB_USER", ""),
			Pass: getEnv("DB_PASS", ""),
			Host: getEnv("DB_HOST", ""),
			Port: getEnv("DB_PORT", ""),
			Name: getEnv("DB_NAME", ""),
		},
		GPT: GPTConfig{
			APIKey: getEnv("GPT_API_KEY", ""),
		},
		Tistory: TistoryConfig{
			Email:    getEnv("TISTORY_EMAIL", ""),
			Password: getEnv("TISTORY_PASSWORD", ""),
			BlogName: getEnv("TISTORY_BLOG_NAME", ""),
			Headless: getEnv("TISTORY_HEADLESS", "true") == "true",
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
