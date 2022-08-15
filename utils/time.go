package utils

import (
	"time"
)

var (
	HOUR  int64 = 3600
	DAY   int64 = HOUR * 24
	MONTH int64 = DAY * 31

	HOURS = []string{
		"0时", "1时", "2时", "3时", "4时", "5时",
		"6时", "7时", "8时", "9时", "10时", "11时",
		"12时", "13时", "14时", "15时", "16时", "17时",
		"18时", "19时", "20时", "21时", "22时", "23时",
	}
	DAYS = []string{
		"1天", "2天", "3天", "4天", "5天", "6天", "7天", "8天", "9天",
		"10天", "11天", "12天", "13天", "14天", "15天", "16天", "17天",
		"18天", "19天", "20天", "21天", "22天", "23天", "24天", "25天",
		"26天", "27天", "28天", "29天", "30天", "31天",
	}
	MONTHS = []string{
		"1月", "2月", "3月", "4月", "5月", "6月",
		"7月", "8月", "9月", "10月", "11月", "12月",
	}
)

func TimeInterval(gap int64) ([]string, []int, int64) {
	if gap > MONTH {
		return MONTHS, make([]int, len(MONTHS)), MONTH
	} else if gap > DAY {
		return DAYS, make([]int, len(DAYS)), DAY
	}
	return HOURS, make([]int, len(HOURS)), HOUR
}

// 获得timestamp当天00:00:00的时间戳，根据当前机器时区计算
func GetZeroClock(timestamp, gap int64) int64 {
	y, m, d := time.Unix(timestamp, 0).Date()
	if gap > MONTH {
		m = time.January
	} else if gap > DAY {
		d = 1
	}
	// fmt.Println(time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix())
	// fmt.Println((timestamp/DAY)*DAY - (8 * HOUR))
	// return (timestamp/DAY)*DAY - (8 * HOUR)	// 也能获取
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix()
}
