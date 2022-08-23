package api

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PerformanceData struct {
	// 首次渲染耗时，时间段内所有渲染时间的平均值
	FirstRunderTime int
	// 首屏时间
	FMPTime ChartData[float32]
	// 首次渲染耗时，时间区间划分
	RunderTime ChartData[float32]
	// 首次可交互时间
	InteractableTime ChartData[float32]
	// dom ready时间
	DomReadyTime ChartData[float32]
	// 页面完全加载时间
	LoadCompleteTime ChartData[float32]
	// 白屏时间
	BlankTime ChartData[float32]
}

func PerformanceTotal(c *gin.Context) {
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
	var data []model.Performance
	searcher.Search(&model.Performance{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "PerformanceTotal: 查询Performance数据失败，该起始时间内没有数据",
			Data: PerformanceData{
				FMPTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				RunderTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				InteractableTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				DomReadyTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				LoadCompleteTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				BlankTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
			},
		})
		return
	}
	// 3. 数据处理
	x, y, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yFMP := make([]float32, len(y))
	yRunder := make([]float32, len(y))
	yIntact := make([]float32, len(y))
	yDom := make([]float32, len(y))
	yLoad := make([]float32, len(y))
	yBlank := make([]float32, len(y))
	var fRunder float32
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		fRunder += data[i].AnalysisTime
		yFMP[count] += data[i].AnalysisTime
		yRunder[count] += data[i].AnalysisTime
		yIntact[count] += data[i].LoadPageTime
		yLoad[count] += data[i].LoadPageTime
		yDom[count] += data[i].DomReadyTime
		yBlank[count] += data[i].BlankTime
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "PerformanceTotal: 查询Performance数据成功",
		Data: PerformanceData{
			FirstRunderTime: int(fRunder) / len(data),
			FMPTime: ChartData[float32]{
				X: x,
				Y: yFMP,
			},
			RunderTime: ChartData[float32]{
				X: x,
				Y: yRunder,
			},
			DomReadyTime: ChartData[float32]{
				X: x,
				Y: yDom,
			},
			InteractableTime: ChartData[float32]{
				X: x,
				Y: yLoad,
			},
			LoadCompleteTime: ChartData[float32]{
				X: x,
				Y: yLoad,
			},
			BlankTime: ChartData[float32]{
				X: x,
				Y: yBlank,
			},
		},
	})
}

func PerformancePage(c *gin.Context) {
	// 1. 解析校验参数
	// 中间件ParseURL已经提前解析过参数了，所以这里的查询和转换并不会出错
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
	var data []model.Performance
	searcher.Search(&model.Performance{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "PerformancePage: 查询Performance数据失败，该起始时间内没有数据",
			Data: PerformanceData{
				FMPTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				RunderTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				InteractableTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				DomReadyTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				LoadCompleteTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
				BlankTime: ChartData[float32]{
					X: []string{},
					Y: []float32{},
				},
			},
		})
		return
	}
	// 3. 数据处理
	x, y, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yFMP := make([]float32, len(y))
	yRunder := make([]float32, len(y))
	yIntact := make([]float32, len(y))
	yDom := make([]float32, len(y))
	yLoad := make([]float32, len(y))
	yBlank := make([]float32, len(y))
	var fRunder float32
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		if data[i].URL == path {
			fRunder += data[i].AnalysisTime
			yFMP[count] += data[i].AnalysisTime
			yRunder[count] += data[i].AnalysisTime
			yIntact[count] += data[i].LoadPageTime
			yLoad[count] += data[i].LoadPageTime
			yDom[count] += data[i].DomReadyTime
			yBlank[count] += data[i].BlankTime
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "PerformancePage: 查询Performance数据成功",
		Data: PerformanceData{
			FirstRunderTime: int(fRunder) / len(data),
			FMPTime: ChartData[float32]{
				X: x,
				Y: yFMP,
			},
			RunderTime: ChartData[float32]{
				X: x,
				Y: yRunder,
			},
			DomReadyTime: ChartData[float32]{
				X: x,
				Y: yDom,
			},
			InteractableTime: ChartData[float32]{
				X: x,
				Y: yLoad,
			},
			LoadCompleteTime: ChartData[float32]{
				X: x,
				Y: yLoad,
			},
			BlankTime: ChartData[float32]{
				X: x,
				Y: yBlank,
			},
		},
	})
}
