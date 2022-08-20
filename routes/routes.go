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
			authed.DELETE("user/item/:projectKey", api.ProjectDelete)
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
		get.GET("access/rank", api.AccessRank)
		get.GET("api/rank", api.ApiRank)

		get.GET("jserror/total", api.JsErrTotal)
		get.GET("access/total", api.AccessTotal)
		get.GET("apierror/total", api.ApiErrTotal)
		get.GET("sourceerror/total", api.SourceErrTotal)
		get.GET("performance/total", api.PerformanceTotal)

		// 该中间件用于检查url中是否携带path参数
		get.Use(middleware.ParseURLMore)

		get.GET("jserror/page", api.JsErrPage)
		get.GET("access/page", api.AccessPage)
		get.GET("apierror/page", api.ApiErrPage)
		get.GET("sourceerror/page", api.SourceErrPage)
		get.GET("performance/page", api.PerformancePage)
	}
	return r
}
