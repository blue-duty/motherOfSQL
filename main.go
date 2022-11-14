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
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
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
			if c.typeName == "int" {
				c.len = 11
			} else {
				c.len = 255
			}
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
	// 2.2 generate sql
	var wg sync.WaitGroup
	for i := 0; i < cfg.Target.Count; i++ {
		wg.Add(1)
		go func() {
			var sqlStr string
			var sqlValue string
			var i int
			sqlStr = fmt.Sprintf("INSERT INTO %s (", cfg.Target.Table)
			cMap.Range(func(key, value interface{}) bool {
				sqlStr += key.(string) + ","
				if value.(column).pri == "PRI" {
					if value.(column).auto == "auto_increment" {
						sqlValue += "NULL,"
					} else {
						if value.(column).typeName == "int" {
							i++
							sqlValue += strconv.Itoa(i) + ","
						} else {
							sqlValue += "'" + utils.UUID() + "',"
						}
					}
				} else {
					switch value.(column).typeName {
					case "int":
						if value.(column).len == 0 {
							sqlValue += strconv.Itoa(utils.GenIntValue(1000)) + ","
						} else {
							sqlValue += strconv.Itoa(utils.GenIntValue(value.(column).len)*10) + ","
						}
					case "char":
						sqlValue += "'" + utils.GenStringValue(value.(column).len) + "',"
					case "varchar":
						sqlValue += "'" + utils.GenStringValue(value.(column).len) + "',"
					case "datetime":
						sqlValue += "'" + utils.GenDatetimeValue() + "',"
					case "float":
						sqlValue += strconv.FormatFloat(utils.GenFloatValue(), 'f', 2, 64) + ","
					case "text":
						sqlValue += "'" + utils.GenChineseValue(10000) + "',"
					}
				}
				return true
			})
			sqlStr = strings.TrimRight(sqlStr, ",") + ") VALUES (" + strings.TrimRight(sqlValue, ",") + ")"
			//fmt.Println(sqlStr)
			_, err = db2.Exec(sqlStr)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
