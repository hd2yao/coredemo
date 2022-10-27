package main

import (
	"goweb/framework"
	"goweb/framework/middleware"
)

/*
	路由规则的需求：
		1.HTTP 方法匹配：RESTful
		2.静态路由匹配
		3.批量通用前缀
		4.动态路由匹配
*/

// 注册路由规则
func registerRouter(core *framework.Core) {
	// 需求1+2：静态路由 + HTTP方法匹配
	core.Get("/user/login", middleware.Test3(), UserLoginController)

	// 需求3：批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		//subjectInnerApi := subjectApi.Group("/info")
		//{
		//	subjectInnerApi.Get("/name", SubjectNameController)
		//}
	}
}
