package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiErrData struct {
	SuccCnt  ChartData[int]
	ErrCnt   ChartData[int]
	SuccRate ChartData[float32]
	// 请求错误数
	TotalErr int
	// 请求错误率
	TotalErrRate float32
	List         []Content
}

type Content struct {
	TimeStamp int64
	Duration  int
	Api       string
	Status    string
	Msg       string
}

func ApiErrTotal(c *gin.Context) {
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
	var data []model.APIError
	searcher.Search(&model.APIError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "ApiErrTotal: 查询ApiErr数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	x, ySuccReq, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yErrReq := make([]int, len(ySuccReq))
	ySuccReqRate := make([]float32, len(ySuccReq))
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	succC, errC := 0, 0
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		if data[i].EventType == "load" {
			ySuccReq[count]++
			succC++
		} else {
			yErrReq[count]++
			errC++
		}
	}
	for i := 0; i < len(ySuccReq); i++ {
		if yErrReq[i] == 0 {
			ySuccReqRate[i] = 1
		} else {
			ySuccReqRate[i] = float32(ySuccReq[i]) / float32(ySuccReq[i]+yErrReq[i])
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询ApiErr数据成功",
		Data: ApiErrData{
			SuccCnt: ChartData[int]{
				X: x,
				Y: ySuccReq,
			},
			ErrCnt: ChartData[int]{
				X: x,
				Y: yErrReq,
			},
			SuccRate: ChartData[float32]{
				X: x,
				Y: ySuccReqRate,
			},
			TotalErr:     errC,
			TotalErrRate: float32(errC) / float32(succC+errC),
		},
	})
}

func ApiErrPage(c *gin.Context) {
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
	var data []model.APIError
	searcher.Search(&model.APIError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "ApiErrPage: 查询ApiErr数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	x, ySuccReq, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yErrReq := make([]int, len(ySuccReq))
	ySuccReqRate := make([]float32, len(ySuccReq))
	// 存放页面的每一条报错信息
	list := make([]Content, 0)
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		if data[i].Pathname == path {
			c := Content{
				TimeStamp: data[i].TimeStamp,
				Duration:  data[i].Duration,
				Api:       data[i].URL,
				Status:    data[i].Status,
				Msg:       data[i].EventType,
			}
			if data[i].EventType == "load" {
				c.Msg = "请求成功"
				ySuccReq[count]++
			} else {
				c.Msg = "请求失败"
				yErrReq[count]++
			}
			list = append(list, c)
		}
	}
	for i := 0; i < len(ySuccReq); i++ {
		if yErrReq[i] == 0 {
			ySuccReqRate[i] = 1
		} else {
			ySuccReqRate[i] = float32(ySuccReq[i]) / float32(ySuccReq[i]+yErrReq[i])
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询ApiErr数据成功",
		Data: ApiErrData{
			SuccCnt: ChartData[int]{
				X: x,
				Y: ySuccReq,
			},
			ErrCnt: ChartData[int]{
				X: x,
				Y: yErrReq,
			},
			SuccRate: ChartData[float32]{
				X: x,
				Y: ySuccReqRate,
			},
			List: list,
		},
	})
}

func ApiRank(c *gin.Context) {
	projectKey := c.Query("projectKey")
	// 不限开始时间
	startTimeStamp := int64(-1)
	endTimeStamp := time.Now().Unix()
	// 2. 查询数据
	var searcher = &service.Searcher{
		ProjectKey:     projectKey,
		StartTimeStamp: startTimeStamp,
		EndTimeStamp:   endTimeStamp,
	}
	var data []model.APIError
	searcher.Search(&model.APIError{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "ApiRank: 查询ApiErr数据失败，该起始时间内没有数据",
		})
		return
	}
	m := make(map[string]int)
	apis := make([]string, 0)
	for i := 0; i < len(data); i++ {
		api := data[i].Pathname
		if _, ok := m[api]; !ok {
			apis = append(apis, api)
		}
		m[api]++
	}
	// 根据访问量从大到小排序
	sort.Slice(apis, func(i, j int) bool {
		return m[apis[i]] > m[apis[j]]
	})
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询AccessRank数据成功",
		Data:   apis,
	})
}
