package config

import (
	"os"

	"github.com/joho/godotenv"
)

// 전체 Config 구조체
type Config struct {
	Server ServerConfig
	DB     DBConfig
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

// 메인 로드 함수
func LoadConfig() *Config {
	// .env 파일 로드
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(".env 파일을 불러올 수 없습니다")
	// }

	if os.Getenv("RENDER") == "" {
		err := godotenv.Load()
		if err != nil {
			// 로그만 출력하고 넘어감
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
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
