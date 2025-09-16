package utils

import (
	"time"
)

func FormatDate(date string) (string, error) {
	// 解析 RFC3339
	startTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}

	// 格式化成 YYYY-MM-DD
	return startTime.Format("2006-01-02"), nil
}

func GetToday() string {
	// 加载新加坡时区
	loc, _ := time.LoadLocation("Asia/Singapore")

	// 当前时间转为新加坡时区
	now := time.Now().In(loc)

	// 格式化成 "2006-01-02"
	return now.Format("2006-01-02")
}
