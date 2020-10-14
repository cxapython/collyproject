package tools

import "time"

// 获取当前时间戳|秒
func GetTimestamp() int64 {
	return time.Now().UnixNano() / 1e9
}

// 获取当前时间戳|毫秒
func GetMicroTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取相对日期[00:00:00]时间戳|秒
func GetZeroTimestamp(years int, months int, days int) int64 {
	timeStr := time.Now().AddDate(years, months, days).Format("2006-01-02 00:00:00")
	timeLocation, _ := time.LoadLocation("Asia/Chongqing")
	timeParse, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, timeLocation)
	return timeParse.Unix()
}