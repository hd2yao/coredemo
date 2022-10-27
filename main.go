package main

import (
	"goweb/framework"
	"goweb/framework/middleware"
	"net/http"
)

func main() {
	core := framework.NewCore()
	// core 中使用 use 注册中间件
	//core.Use(
	//	middleware.Test1(),
	//	middleware.Test2())

	// 为所有的路由都设置 Recovery 中间件
	core.Use(middleware.Recovery())
	//core.Use(middleware.Cost())
	// core.Use(middleware.Timeout(1 * time.Second))

	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
