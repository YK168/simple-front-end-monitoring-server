package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlankData struct {
	BlankRate int
	Data      ChartData[int]
}

func BlankPage(c *gin.Context) {
	// 1. 解析校验参数
	// 中间件Parse已经提前解析过参数了，所以这里的查询和转换并不会出错
	projectKey := c.Query("projectKey")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	path := c.Query("path")
	startTimeStamp, _ := strconv.ParseInt(startTime, 10, 64)
	endTimeStamp, _ := strconv.ParseInt(endTime, 10, 64)
	// 2. 查询数据
	var searcher = &service.Searcher{
		ProjectKey:     projectKey,
		StartTimeStamp: startTimeStamp,
		EndTimeStamp:   endTimeStamp,
	}
	var data []model.BlankError
	searcher.Search(&model.BlankError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "BlankPage: 查询Blank Error数据失败，该起始时间内没有数据",
		})
		return
	}
	x, y, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	totalBlank := 0
	tBlank := 0
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			tBlank = totalBlank
			count++
		}
		if data[i].URL == path {
			totalBlank++
			// 记录这个区间中Pv和Uv的数据
			y[count] = totalBlank - tBlank
		}
	}
	var data2 []model.Access
	searcher.Search(&model.Access{}, &data2)
	pv := 0
	for i := 0; i < len(data2); i++ {
		if data2[i].URL == path {
			pv++
		}
	}
	if pv == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "BlankPage: 该起始时间内没有" + path + "的Access记录，无法计算白屏率",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询BlankPage数据成功",
		Data: BlankData{
			BlankRate: totalBlank / pv,
			Data: ChartData[int]{
				X: x,
				Y: y,
			},
		},
	})
}
