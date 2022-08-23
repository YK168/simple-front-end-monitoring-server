package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JsErrData struct {
	TotalErr int
	Data     ChartData[int]
}

type SourceErrData struct {
	Data ChartData[int]
	// 请求错误数
	TotalErr int
}

type ChartData[T int | float32] struct {
	X []string
	Y []T
}

func JsErrTotal(c *gin.Context) {
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
	var data []model.JSError
	searcher.Search(&model.JSError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "查询JsError数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	x, y, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data) && count < len(x); i++ {
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		y[count]++
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询JsError数据成功",
		Data: JsErrData{
			TotalErr: len(data),
			Data: ChartData[int]{
				X: x,
				Y: y,
			},
		},
	})
}

func JsErrPage(c *gin.Context) {
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
	var data []model.JSError
	searcher.Search(&model.JSError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "查询JsError数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	type Data struct {
		Msg      string
		Time     string
		Position string
	}
	d := make([]Data, 0, len(data))
	for i := 0; i < len(data); i++ {
		if data[i].URL == path {
			d = append(d, Data{
				Msg:      data[i].Message,
				Time:     utils.TimeStampToDate(data[i].TimeStamp),
				Position: data[i].Position,
			})
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询JsError数据成功",
		Data:   d,
	})
}
