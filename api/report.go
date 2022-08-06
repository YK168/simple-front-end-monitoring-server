package api

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"

	"github.com/gin-gonic/gin"
)

// JS error 上报
func JsErrReport(c *gin.Context) {
	var jsErr service.JsErrorService
	if err := c.ShouldBind(&jsErr); err == nil {
		res := jsErr.Report()
		c.JSON(res.Status, res)
	} else {
		log.Println("解析json数据失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "解析json数据失败",
			Error:  err.Error(),
		})
	}
}

// Api error 上报
func ApiErrReport(c *gin.Context) {
	var apiErr service.ApiErrorService
	if err := c.ShouldBind(&apiErr); err == nil {
		res := apiErr.Report()
		c.JSON(res.Status, res)
	} else {
		log.Println("解析json数据失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "解析json数据失败",
			Error:  err.Error(),
		})
	}
}

// 资源 error 上报
func SourceErrReport(c *gin.Context) {
	var sourceErr service.SourceErrorService
	if err := c.ShouldBind(&sourceErr); err == nil {
		res := sourceErr.Report()
		c.JSON(res.Status, res)
	} else {
		log.Println("解析json数据失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "解析json数据失败",
			Error:  err.Error(),
		})
	}
}

// 性能数据 上报
func PerformanceReport(c *gin.Context) {
	var performance service.PerformanceService
	if err := c.ShouldBind(&performance); err == nil {
		res := performance.Report()
		c.JSON(res.Status, res)
	} else {
		log.Println("解析json数据失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "解析json数据失败",
			Error:  err.Error(),
		})
	}
}
