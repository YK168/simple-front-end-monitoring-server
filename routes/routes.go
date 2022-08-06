package routes

import (
	"simple_front_end_monitoring_server/api"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/v1")
	{
		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		// TODO: 创建项目
	}
	reporter := r.Group("api/reporter")
	{
		// 数据上报
		reporter.POST("jserror", api.JsErrReport)
		reporter.POST("apierror", api.ApiErrReport)
		reporter.POST("sourceerror", api.SourceErrReport)
		reporter.POST("performance", api.PerformanceReport)
	}
	return r
}
