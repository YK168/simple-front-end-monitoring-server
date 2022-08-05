package api

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/service"

	"github.com/gin-gonic/gin"
)

// JS error 上报
func JsErrReport(c *gin.Context) {
	var jsErr service.JsErrorService
	if err := c.ShouldBind(&jsErr); err == nil {
		log.Println(err)
		res := jsErr.Report()
		c.JSON(res.Status, res)
	} else {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
	}
}

// Api error 上报
func ApiErrReport(c *gin.Context) {
	var apiErr service.ApiErrorService
	if err := c.ShouldBind(&apiErr); err == nil {
		res := apiErr.Report()
		c.JSON(res.Status, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
