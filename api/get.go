package api

import (
	"log"
	"net/http"
	"net/url"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccessData struct {
	TotalPV int
	TotalUV int
	PVData  ChartData[int]
	UVData  ChartData[int]
}

type JsErrData struct {
	TotalErr int
	Data     ChartData[int]
}

type ApiErrData struct {
	SuccCnt  ChartData[int]
	ErrCnt   ChartData[int]
	SuccRate ChartData[float32]
	// 请求错误数
	TotalErr int
	// 请求错误率
	TotalErrRate float32
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
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, y, gap := utils.TimeInterval(gap)
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data); i++ {
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		y[count]++
	}
	l, r := utils.GetBorder(y)
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询JsError数据成功",
		Data: JsErrData{
			TotalErr: len(data),
			Data: ChartData[int]{
				X: x[l:r],
				Y: y[l:r],
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
		u, err := url.Parse(data[i].URL)
		if err != nil {
			log.Println("AccessPage: 解析出错", err)
			continue
		}
		if u.Path == path {
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
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "TotleAccessGet: 查询Access数据失败，该起始时间内没有数据",
		})
		return
	}
	// 3. 数据处理
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, yPV, gap := utils.TimeInterval(gap)
	yUV := make([]int, len(yPV))
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	totalPV, totalUV := 0, 0
	tPV, tUV := 0, 0
	hasUser := map[string]bool{}
	for i := 0; i < len(data); i++ {
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
	l1, r1 := utils.GetBorder(yPV)
	l2, r2 := utils.GetBorder(yUV)
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询Access数据成功",
		Data: AccessData{
			TotalPV: totalPV,
			TotalUV: totalUV,
			PVData: ChartData[int]{
				X: x[l1:r1],
				Y: yPV[l1:r1],
			},
			UVData: ChartData[int]{
				X: x[l2:r2],
				Y: yUV[l2:r2],
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
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "AccessPage: 查询Access数据失败，该起始时间内没有数据",
		})
		return
	}
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, yPV, gap := utils.TimeInterval(gap)
	yUV := make([]int, len(yPV))
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	totalPV, totalUV := 0, 0
	tPV, tUV := 0, 0
	hasUser := map[string]bool{}
	for i := 0; i < len(data); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			tPV = totalPV
			tUV = totalUV
			count++
		}
		u, err := url.Parse(data[i].URL)
		if err != nil {
			log.Println("AccessPage: 解析出错", err)
			continue
		}
		if u.Path == path {
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
	l1, r1 := utils.GetBorder(yPV)
	l2, r2 := utils.GetBorder(yUV)
	c.JSON(http.StatusOK, utils.Response{
		Status: http.StatusOK,
		Msg:    "查询AccessPage数据成功",
		Data: AccessData{
			TotalPV: totalPV,
			TotalUV: totalUV,
			PVData: ChartData[int]{
				X: x[l1:r1],
				Y: yPV[l1:r1],
			},
			UVData: ChartData[int]{
				X: x[l2:r2],
				Y: yUV[l2:r2],
			},
		},
	})
}

func AccessRank(c *gin.Context) {
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
			Msg:    "AccessRank: 查询Access数据失败，该起始时间内没有数据",
		})
		return
	}
	m := make(map[string]int)
	paths := make([]string, 0)
	for i := 0; i < len(data); i++ {
		u, err := url.Parse(data[i].URL)
		if err != nil {
			log.Println("AccessRank: 解析出错", err)
			continue
		}
		path := u.Path
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
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, ySuccReq, gap := utils.TimeInterval(gap)
	yErrReq := make([]int, len(ySuccReq))
	ySuccReqRate := make([]float32, len(ySuccReq))
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	succC, errC := 0, 0
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data); i++ {
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
	l, r := utils.Get2Border(ySuccReq, yErrReq)
	x = x[l:r]
	ySuccReq = ySuccReq[l:r]
	yErrReq = yErrReq[l:r]
	ySuccReqRate = ySuccReqRate[l:r]
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
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, ySuccReq, gap := utils.TimeInterval(gap)
	yErrReq := make([]int, len(ySuccReq))
	ySuccReqRate := make([]float32, len(ySuccReq))
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	for i := 0; i < len(data); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		u, err := url.Parse(data[i].Pathname)
		if err != nil {
			log.Println("ApiErrPage: 解析出错", err)
			continue
		}
		if u.Hostname() == path {
			if data[i].EventType == "load" {
				ySuccReq[count]++
			} else {
				yErrReq[count]++
			}
		}
	}
	l, r := utils.Get2Border(ySuccReq, yErrReq)
	x = x[l:r]
	ySuccReq = ySuccReq[l:r]
	yErrReq = yErrReq[l:r]
	ySuccReqRate = ySuccReqRate[l:r]
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
		},
	})
}

func ApiRank(c *gin.Context) {
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
			Msg:    "ApiRank: 查询ApiErr数据失败，该起始时间内没有数据",
		})
		return
	}
	m := make(map[string]int)
	apis := make([]string, 0)
	for i := 0; i < len(data); i++ {
		u, err := url.Parse(data[i].Pathname)
		if err != nil {
			log.Println("AccessRank: 解析出错", err)
			continue
		}
		host := u.Hostname()
		if _, ok := m[host]; !ok {
			apis = append(apis, host)
		}
		m[host]++
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
	gap := searcher.EndTimeStamp - searcher.StartTimeStamp
	x, y, gap := utils.TimeInterval(gap)
	var t1 int64 = utils.GetZeroClock(data[0].TimeStamp, gap)
	var t2 int64 = t1 + gap
	count := 0 // 计算在第几个时间区间
	errC := 0
	for i := 0; i < len(data); i++ {
		// 将t2移动到有数据的区间
		for data[i].TimeStamp > t2 {
			t2 += gap
			count++
		}
		u, err := url.Parse(data[i].URL)
		if err != nil {
			log.Println("SourceErrPage: 解析出错", err)
			continue
		}
		if u.Path == path {
			errC++
			y[count]++
		}
	}
	l, r := utils.GetBorder(y)
	x = x[l:r]
	y = y[l:r]
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
