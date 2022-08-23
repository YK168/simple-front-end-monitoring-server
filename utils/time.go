package utils

import (
	"time"

	"github.com/golang-module/carbon/v2"
)

var (
	HOUR  int64 = 3600
	DAY   int64 = HOUR * 24
	MONTH int64 = DAY * 31
)

// 判断是月区间还是天区间还是时区间
// BUG: 当月天数少于31天时，无法正确判断区间
func TimeInterval(startTime, endTime int64) ([]string, []int, int64) {
	x := make([]string, 0)
	gap := endTime - startTime
	format := ""
	if gap > MONTH {
		gap = MONTH
		format = "2006-01"
	} else if gap > DAY {
		gap = DAY
		format = "2006-01-02"
	} else {
		gap = HOUR
		format = "15:04:05"
	}
	origin := startTime
	j := 1
	for startTime <= endTime {
		x = append(x, time.Unix(startTime, 0).Format(format))
		if gap == MONTH {
			startTime = carbon.CreateFromTimestamp(origin).StartOfMonth().AddMonths(j).Timestamp()
			j++
		} else {
			startTime += gap
		}
	}
	return x, make([]int, len(x)), gap
}

// 根据gap获得某天00:00:00的时间戳，
// 或某个月第一天的时间戳，
// 或某年第一个月的时间戳
// 是根据当前机器时区计算的
func GetZeroClock(timestamp, gap int64) int64 {
	y, m, d := time.Unix(timestamp, 0).Date()
	h := 0
	if gap > MONTH {
		m = time.January
	} else if gap > DAY {
		d = 1
	} else {
		h, _, _ = time.Unix(timestamp, 0).Clock()
	}
	// return (timestamp/DAY)*DAY - (8 * HOUR)	// 也能获取
	return time.Date(y, m, d, h, 0, 0, 0, time.Local).Unix()
}

func TimeStampToDate(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02")
}
