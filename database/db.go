package db

import (
    "database/sql"
    "time"

    "fmt"

    "github.com/go-sql-driver/mysql"
    _ "github.com/go-sql-driver/mysql"
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
    cfg := mysql.Config{
        User:                 conf.User,
        Passwd:               conf.Pw,
        Net:                  "tcp",
        Addr:                 conf.Address + ":" + conf.Port,
        DBName:               conf.Scheme,
        AllowNativePasswords: true,
        ParseTime:            true,
    }
    ret, err := sql.Open("mysql", cfg.FormatDSN())
    // what to do?
    if err != nil {
        return nil, err
    }

    ret.SetMaxIdleConns(conf.MaxIdle)
    ret.SetMaxOpenConns(conf.MaxOpen)
    ret.SetConnMaxLifetime(time.Minute * time.Duration(conf.MaxLifeTime))

    if err = ret.Ping(); err != nil {
        return nil, err
    }

    DB = ret

    return ret, nil
}

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
        me, ok := err.(*mysql.MySQLError)
        if !ok {
            return -1, err
        }
        if me.Number == 1062 {
            return 0, err
        } else {
            return -1, err
        }
    }

    nRow, err := ret.RowsAffected()
    if err != nil {
        return 01, err
    }

    return int(nRow), nil
}
