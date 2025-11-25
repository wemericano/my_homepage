package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// DB database global
var DB *sql.DB

type DBConfig struct {
	Address     string
	Port        string
	User        string
	Pw          string
	Scheme      string
	MaxIdle     int
	MaxOpen     int
	MaxLifeTime int
}

func Open(conf DBConfig) (*sql.DB, error) {
	log.Println("[DB] SQL Server 연결 설정 중...")

	connStr := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s",
		conf.Address, conf.User, conf.Pw, conf.Port, conf.Scheme,
	)
	log.Printf("[DB] 연결 문자열 생성 완료 (주소: %s:%s, DB: %s)", conf.Address, conf.Port, conf.Scheme)

	ret, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Printf("[DB] sql.Open 실패: %v", err)
		return nil, fmt.Errorf("sql.Open 실패: %w", err)
	}
	log.Println("[DB] sql.Open 성공")

	ret.SetMaxIdleConns(conf.MaxIdle)
	ret.SetMaxOpenConns(conf.MaxOpen)
	ret.SetConnMaxLifetime(time.Minute * time.Duration(conf.MaxLifeTime))

	log.Println("[DB] Ping 시도 중...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ret.PingContext(ctx); err != nil {
		log.Printf("[DB] Ping 실패: %v", err)
		log.Printf("[DB] 연결 정보 확인 필요 - Host: %s, Port: %s, User: %s, DB: %s",
			conf.Address, conf.Port, conf.User, conf.Scheme)
		return nil, fmt.Errorf("ping 실패: %w", err)
	}
	log.Println("[DB] Ping 성공 - DB 연결 완료")

	DB = ret
	return ret, nil
}

// func Open(conf DBConfig) (*sql.DB, error) {
// 	log.Println("[DB] MySQL 연결 설정 중...")
//
// 	// 포트 확인 (1433은 SQL Server 포트, MySQL은 보통 3306)
// 	if conf.Port == "1433" {
// 		log.Println("[DB] 경고: Port 1433은 SQL Server 기본 포트입니다. MySQL은 보통 3306을 사용합니다.")
// 	}
//
// 	cfg := mysql.Config{
// 		User:                 conf.User,
// 		Passwd:               conf.Pw,
// 		Net:                  "tcp",
// 		Addr:                 conf.Address + ":" + conf.Port,
// 		DBName:               conf.Scheme,
// 		AllowNativePasswords: true,
// 		ParseTime:            true,
// 		Timeout:              10 * time.Second, // 연결 타임아웃
// 		ReadTimeout:          30 * time.Second, // 읽기 타임아웃
// 		WriteTimeout:         30 * time.Second, // 쓰기 타임아웃
// 	}
//
// 	dsn := cfg.FormatDSN()
// 	log.Printf("[DB] DSN 생성 완료 (주소: %s:%s, DB: %s)", conf.Address, conf.Port, conf.Scheme)
//
// 	ret, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Printf("[DB] sql.Open 실패: %v", err)
// 		return nil, fmt.Errorf("sql.Open 실패: %w", err)
// 	}
// 	log.Println("[DB] sql.Open 성공")
//
// 	ret.SetMaxIdleConns(conf.MaxIdle)
// 	ret.SetMaxOpenConns(conf.MaxOpen)
// 	ret.SetConnMaxLifetime(time.Minute * time.Duration(conf.MaxLifeTime))
//
// 	log.Println("[DB] Ping 시도 중...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	if err = ret.PingContext(ctx); err != nil {
// 		log.Printf("[DB] Ping 실패: %v", err)
// 		log.Printf("[DB] 연결 정보 확인 필요 - Host: %s, Port: %s, User: %s, DB: %s",
// 			conf.Address, conf.Port, conf.User, conf.Scheme)
// 		return nil, fmt.Errorf("Ping 실패: %w", err)
// 	}
// 	log.Println("[DB] Ping 성공 - DB 연결 완료")
//
// 	DB = ret
//
// 	return ret, nil
// }

// TODO
func Reopen(db *sql.DB) (*sql.DB, error) {
	return db, nil
}

func TableData2Map(rows *sql.Rows) []map[string]interface{} {
	cols_nm, _ := rows.Columns()
	cols_cnt := len(cols_nm)
	values := make([]interface{}, cols_cnt)
	value_ptrs := make([]interface{}, cols_cnt)

	var ret []map[string]interface{}
	for rows.Next() {

		for i := range cols_nm {
			value_ptrs[i] = &values[i]
		}

		rows.Scan(value_ptrs...)

		var r_data = make(map[string]interface{})
		for i, col_nm := range cols_nm {
			v := values[i]
			b, e := v.([]byte)
			var val interface{}
			if e {
				val = string(b)
			} else {
				val = v
			}

			if val == nil {
				r_data[col_nm] = ""
			} else {
				r_data[col_nm] = fmt.Sprint(val)
			}

			//swlog.Dev(col_nm + ":" + val.(string))
		}
		ret = append(ret, r_data)
	}
	return ret
}

func Mysql_fetch_rows(db *sql.DB, sqlquery string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	retdata := TableData2Map(rows)

	return retdata, nil
}

func Mysql_insert_rows_one(db *sql.DB, sqlquery string) (int, error) {
	ret, err := db.Exec(sqlquery)
	if err != nil {
		// SQL Server는 MySQL과 다른 에러 타입을 사용합니다
		// 중복 키 에러 처리는 필요시 추가
		return -1, err
	}

	nRow, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(nRow), nil
}
