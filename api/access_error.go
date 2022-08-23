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

type AccessData struct {
	TotalPV int
	TotalUV int
	PVData  ChartData[int]
	UVData  ChartData[int]
}

func AccessTotal(c *gin.Context) {
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
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "TotleAccessGet: 查询Access数据失败，该起始时间内没有数据",
			Data: AccessData{
				PVData: ChartData[int]{
					X: []string{},
					Y: []int{},
				},
				UVData: ChartData[int]{
					X: []string{},
					Y: []int{},
				},
			},
		})
		return
	}
	// 3. 数据处理
	x, yPV, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yUV := make([]int, len(yPV))
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	totalPV, totalUV := 0, 0
	tPV, tUV := 0, 0
	hasUser := map[string]bool{}
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			tPV = totalPV
			tUV = totalUV
			count++
		}
		totalPV++
		// 根据Cookie来区分不同用户，没有Cookie的记录认为是两个不同用户
		if data[i].Cookie == "" || !hasUser[data[i].Cookie] {
			totalUV++
			hasUser[data[i].Cookie] = true
		}
		// 记录这个区间中Pv和Uv的数据
		yPV[count] = totalPV - tPV
		yUV[count] = totalUV - tUV
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询Access数据成功",
		Data: AccessData{
			TotalPV: totalPV,
			TotalUV: totalUV,
			PVData: ChartData[int]{
				X: x,
				Y: yPV,
			},
			UVData: ChartData[int]{
				X: x,
				Y: yUV,
			},
		},
	})
}

func AccessPage(c *gin.Context) {
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
	var data []model.Access
	searcher.Search(&model.Access{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "AccessPage: 查询Access数据失败，该起始时间内没有数据",
			Data: AccessData{
				PVData: ChartData[int]{
					X: []string{},
					Y: []int{},
				},
				UVData: ChartData[int]{
					X: []string{},
					Y: []int{},
				},
			},
		})
		return
	}
	// gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, yPV, gap := utils.TimeInterval(searcher.StartTimeStamp, searcher.EndTimeStamp)
	yUV := make([]int, len(yPV))
	var t1 int64 = utils.GetZeroClock(searcher.StartTimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	totalPV, totalUV := 0, 0
	tPV, tUV := 0, 0
	hasUser := map[string]bool{}
	for i := 0; i < len(data) && count < len(x); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			tPV = totalPV
			tUV = totalUV
			count++
		}
		if data[i].URL == path {
			totalPV++
			// 根据Cookie来区分不同用户，没有Cookie的记录认为是两个不同用户
			if data[i].Cookie == "" || !hasUser[data[i].Cookie] {
				totalUV++
				hasUser[data[i].Cookie] = true
			}
			// 记录这个区间中Pv和Uv的数据
			yPV[count] = totalPV - tPV
			yUV[count] = totalUV - tUV
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询AccessPage数据成功",
		Data: AccessData{
			TotalPV: totalPV,
			TotalUV: totalUV,
			PVData: ChartData[int]{
				X: x,
				Y: yPV,
			},
			UVData: ChartData[int]{
				X: x,
				Y: yUV,
			},
		},
	})
}

func AccessRank(c *gin.Context) {
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
	var data []model.Access
	searcher.Search(&model.Access{}, &data)
	if len(data) == 0 {
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "AccessRank: 查询Access数据失败，该起始时间内没有数据",
			Data:   []string{},
		})
		return
	}
	m := make(map[string]int)
	paths := make([]string, 0)
	for i := 0; i < len(data); i++ {
		path := data[i].URL
		if _, ok := m[path]; !ok {
			paths = append(paths, path)
		}
		m[path]++
	}
	// 根据访问量从大到小排序
	sort.Slice(paths, func(i, j int) bool {
		return m[paths[i]] > m[paths[j]]
	})
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询AccessRank数据成功",
		Data:   paths,
	})
}
