package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func JsErrGet(c *gin.Context) {
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
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	// 3.1 先判断是月区间还是天区间还是时区间
	x, y, gap := utils.TimeInterval(gap)
	// 借助标准库，做到自由获取第几年第几个月第几天的时间戳
	// 这样可以获得某天00:00:00的时间戳，也可以获得某个月第一天的时间戳，也可以获得某年第一个月的时间戳
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	// 3.2将查询出来的数据填充x和y数组
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data); i++ {
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		// 根据逻辑修改，jsErr只计算访问数，所以++就好
		y[count]++
	}
	l, r := utils.GetBorder(y)
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询JsError数据成功",
		Data:   []any{x[l:r], y[l:r]},
	})
}

func AccessGet(c *gin.Context) {
	// 1. 解析校验参数
	// 中间件ParseURL已经提前解析过参数了，所以这里的查询和转换并不会出错
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
	var data []model.Access
	searcher.Search(&model.Access{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "查询Access数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	// 3.1 先判断是月区间还是天区间还是时区间
	x, y, gap := utils.TimeInterval(gap)
	// 借助标准库，做到自由获取第几年第几个月第几天的时间戳
	// 这样可以获得某天00:00:00的时间戳，也可以获得某个月第一天的时间戳，也可以获得某年第一个月的时间戳
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	// 3.2将查询出来的数据填充x和y数组
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data); i++ {
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		// 根据逻辑修改，access只计算访问数，所以++就好
		y[count]++
	}
	l, r := utils.GetBorder(y)
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询Access数据成功",
		Data:   []any{x[l:r], y[l:r]},
	})
}
