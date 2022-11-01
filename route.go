package main

import (
	"github.com/haidongXX/coredemo/framework/gin"
	"github.com/haidongXX/coredemo/framework/middleware"
)

/*
	路由规则的需求：
		1.HTTP 方法匹配：RESTful
		2.静态路由匹配
		3.批量通用前缀
		4.动态路由匹配
*/

// 注册路由规则
func registerRouter(core *gin.Engine) {
	// 需求1+2：静态路由 + HTTP方法匹配
	core.GET("/user/login", middleware.Test3(), UserLoginController)

	// 需求3：批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 动态路由
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)

		//subjectInnerApi := subjectApi.Group("/info")
		//{
		//	subjectInnerApi.Get("/name", SubjectNameController)
		//}
	}
}
