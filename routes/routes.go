package routes

import (
	"simple_front_end_monitoring_server/api"
	"simple_front_end_monitoring_server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	// 设置跨域
	r.Use(cors.Default())
	// 用户操作
	v1 := r.Group("api/v1")
	{
		// 注册登录
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		authed := v1.Group("/") //需要登陆保护
		// 设置JWT中间件
		authed.Use(middleware.JWT)
		{
			// 创建项目，返回生成的project key
			authed.POST("user/item", api.ProjectCreate)
			// 删除项目
			authed.DELETE("user/item", api.ProjectDelete)
			// 修改项目，现在projectKey是根据md5(number + title)生成的
			// 如果更新title，则projectKey需要重新生成
			// 如果只更新title而不更新projectKet，会导致后续再有名叫title的项目生成时
			// 新旧title的数据会混淆
			// 涉及到projectKey和关联的监控数据，更新很复杂，暂时不实现
			// authed.PUT("user/item", api.ProjectUpdate)
			// 查询项目
			authed.GET("user/item", api.ProjectSearch)
		}
	}
	reporter := r.Group("api/reporter")
	{
		// 数据上报
		reporter.POST("jserror", api.JsErrReport)
		reporter.POST("apierror", api.ApiErrReport)
		reporter.POST("sourceerror", api.SourceErrReport)
		reporter.POST("performance", api.PerformanceReport)
		reporter.POST("access", api.AccessReport)
	}
	get := r.Group("api/get")
	{
		get.Use(middleware.ParseURL)

		// 数据请求，返回用于echarts生成图表的x轴和y轴数组
		get.GET("jserror", api.JsErrGet)
		// get.GET("apierror", api.ApiErrGet)
		// get.GET("sourceerror", api.SourceErrGet)
		// get.GET("performance", api.PerformanceGet)
		get.GET("totalaccess", api.TotleAccessGet)
	}
	return r
}
