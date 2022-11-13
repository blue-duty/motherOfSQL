package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"motherOfSQL/config"
	"motherOfSQL/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type column struct {
	typeName string
	auto     string
	pri      string
	len      int
}

func main() {
	t1 := time.Now()
	var cMap sync.Map
	cfg, err := config.ReadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}
	db1, err := sql.Open(cfg.DB.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, "information_schema"))
	if err != nil {
		panic(err)
	}
	defer func(db1 *sql.DB) {
		_ = db1.Close()
	}(db1)
	//1. get column type and name form table
	rows, err := db1.Query("SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY,EXTRA,COLUMN_TYPE FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", cfg.Target.Schema, cfg.Target.Table)
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var c column
		var t string
		var name string
		err = rows.Scan(&name, &c.typeName, &c.pri, &c.auto, &t)
		if err != nil {
			panic(err)
		}
		if strings.Contains(t, "(") {
			c.len, _ = strconv.Atoi(utils.TakeParentheses(t))
		} else {
			c.len = 0
		}
		cMap.Store(name, c)
	}

	// 2. generate 100 sql statements with value by cMap
	db2, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.Target.Schema))
	if err != nil {
		panic(err)
	}
	defer func(db2 *sql.DB) {
		_ = db2.Close()
	}(db2)
	// 2.1 get column name
	var columns []string
	var values []interface{}
	cMap.Range(func(key, value interface{}) bool {
		columns = append(columns, key.(string))
		values = append(values, value.(column))
		return true
	})
	fmt.Println(values)
	// 2.2 generate sql
	var m sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < cfg.Target.Count; i++ {
		wg.Add(1)
		go func() {
			var sqlStr string
			sqlStr = fmt.Sprintf("INSERT INTO %s (", cfg.Target.Table)
			m.Lock()
			for _, v := range columns {
				sqlStr += v + ","
			}
			sqlStr = sqlStr[:len(sqlStr)-1] + ") VALUES ("
			for i, v := range values {
				if v.(column).pri == "PRI" {
					if v.(column).auto == "auto_increment" {
						sqlStr += "NULL,"
					} else {
						if v.(column).typeName == "int" {
							sqlStr += strconv.Itoa(i+1) + ","
						} else {
							sqlStr += "'" + utils.UUID() + "',"
						}
					}
				} else {
					switch v.(column).typeName {
					case "int":
						if v.(column).len == 0 {
							sqlStr += strconv.Itoa(utils.GenIntValue(1000)) + ","
						} else {
							sqlStr += strconv.Itoa(utils.GenIntValue(v.(column).len)*10) + ","
						}
					case "char":
						if v.(column).len == 0 {
							sqlStr += "'" + utils.GenStringValue(200) + "',"
						} else {
							sqlStr += "'" + utils.GenStringValue(v.(column).len) + "',"
						}
					case "varchar":
						if v.(column).len == 0 {
							sqlStr += "'" + utils.GenStringValue(200) + "',"
						} else {
							sqlStr += "'" + utils.GenStringValue(v.(column).len) + "',"
						}
					case "datetime":
						sqlStr += "'" + utils.GenDatetimeValue() + "',"
					case "float":
						sqlStr += strconv.FormatFloat(utils.GenFloatValue(), 'f', 2, 64) + ","
					case "text":
						sqlStr += "'" + utils.GenChineseValue(10000) + "',"
					}
				}

			}
			m.Unlock()
			sqlStr = sqlStr[:len(sqlStr)-1] + ")"
			fmt.Println(sqlStr)
			_, err = db2.Exec(sqlStr)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(t1).String())
}
