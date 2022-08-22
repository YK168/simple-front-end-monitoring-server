package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SourceErrTotal(c *gin.Context) {
	// 1. 解析校验参数
	// 中间件Parse已经提前解析过参数了，所以这里的查询和转换并不会出错
	projectKey := c.Query("projectKey")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	startTimeStamp, _ := strconv.ParseInt(startTime, 10, 64)
	endTimeStamp, _ := strconv.ParseInt(endTime, 10, 64)
	// 2. 查询数据
	var searcher = &service.Searcher{
		ProjectKey:     projectKey,
		StartTimeStamp: startTimeStamp,
		EndTimeStamp:   endTimeStamp,
	}
	var data []model.SourceError
	searcher.Search(&model.SourceError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "SourceErrTotal: 查询SourceErr数据失败，该起始时间内没有数据",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "SourceErrTotal: 查询SourceErr数据成功",
		Data: SourceErrData{
			TotalErr: len(data),
		},
	})
}

func SourceErrPage(c *gin.Context) {
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
	var data []model.SourceError
	searcher.Search(&model.SourceError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "SourceErrPage: 查询SourceErr数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	x, y, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	errC := 0
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		if data[i].URL == path {
			errC++
			y[count]++
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "SourceErrPage: 查询SourceErr数据成功",
		Data: SourceErrData{
			Data: ChartData[int]{
				X: x,
				Y: y,
			},
			TotalErr: errC,
		},
	})
}
