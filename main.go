package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

// 生成指定时间范围内的随机时间
func randomTimeWithinRange(start, end string) string {
	layout := "15:04:05"
	startTime, _ := time.Parse(layout, start)
	endTime, _ := time.Parse(layout, end)

	delta := endTime.Sub(startTime).Seconds()
	randomSeconds := rand.Int63n(int64(delta))
	randomTime := startTime.Add(time.Duration(randomSeconds) * time.Second)

	return randomTime.Format(layout)
}

// 插入新的记录到数据库
func insertNewRecord(dbPath, inputDate string) {
	// 连接到 SQLite 数据库
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 定义变量存储上一次记录的信息
	var (
		lastID, newFlag                                         int
		newStratpoint, newEndpoint, newPathline, newNum, newTmd sql.NullString
	)

	// 查询最后一条记录
	err = db.QueryRow("SELECT id, stratpoint, endpoint, pathline, num, flag, tmd FROM record ORDER BY id DESC LIMIT 1;").Scan(&lastID, &newStratpoint, &newEndpoint, &newPathline, &newNum, &newFlag, &newTmd)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Failed to fetch last record: %v", err)
	}

	// 设置新记录的 ID 和默认值
	newID := lastID + 1
	if err == sql.ErrNoRows {
		newID = 1
		newStratpoint = sql.NullString{String: "InitialStartPoint", Valid: true}
		newEndpoint = sql.NullString{String: "InitialEndPoint", Valid: true}
		newPathline = sql.NullString{String: "InitialPathLine", Valid: true}
		newNum = sql.NullString{String: "InitialNum", Valid: true}
		newFlag = 0
		newTmd = sql.NullString{String: "InitialTmd", Valid: true}
	}

	// 随机生成记录的其他字段值
	distance := fmt.Sprintf("%.4f", 5001+rand.Float64()*(5023-5001))
	duration := rand.Intn(1489-1350+1) + 1350
	averageSpeed := fmt.Sprintf("%.1f", 3002+rand.Float64()*(3356-3002))
	step := rand.Intn(3325-3145+1) + 3145
	calorie := rand.Intn(302-274+1) + 274

	// 生成随机时间并组合成完整的日期时间字符串
	randomTime := randomTimeWithinRange("17:30:00", "21:00:00")
	inputDatetime := fmt.Sprintf("%s %s", inputDate, randomTime)

	// 插入新记录到数据库
	_, err = db.Exec(`
		INSERT INTO record (id, stratpoint, endpoint, pathline, distance, duration, averagespeed, date, num, step, calorie, flag, tmd)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, newID, newStratpoint.String, newEndpoint.String, newPathline.String, distance, duration, averageSpeed, inputDatetime, newNum.String, step, calorie, newFlag, newTmd.String)
	if err != nil {
		log.Fatalf("Failed to insert new record: %v", err)
	}

	fmt.Printf("New record inserted with ID: %d\n", newID)
}

func main() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 数据库路径和输入日期
	dbPath := "record.db"
	inputDate := "2024-10-30"

	// 调用插入函数
	insertNewRecord(dbPath, inputDate)
}
