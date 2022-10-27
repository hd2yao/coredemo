package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构
type Core struct {
	router map[string]*Tree
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = newTree()
	router["POST"] = newTree()
	router["PUT"] = newTree()
	router["DELETE"] = newTree()
	return &Core{router: router}
}

// 注册路由

// Get 匹配 GET 方法，增加路由规则
func (c *Core) Get(url string, handler ControllerHandler) {
	if err := c.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Post 匹配 POST 方法
func (c *Core) Post(url string, handler ControllerHandler) {
	if err := c.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put 匹配 PUT 方法
func (c *Core) Put(url string, handler ControllerHandler) {
	if err := c.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete 匹配 DELETE 方法
func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Group 通用前缀
//func (c *Core) Group(prefix string) IGroup {
//	return NewGroup(c, prefix)
//}

// FindRouteByRequest 匹配路由，如果没有匹配到，返回 nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	// uri 和 method 全部转换成大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层 map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// 框架核心结构实现 Handler 接口
// 所有请求都进入这个函数，这个函数负责路由分发
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// 封装自定义 context
	ctx := NewContext(request, response)

	// 寻找路由
	router := c.FindRouteByRequest(request)
	if router == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}

	// 调用路由函数，如果返回 err 代表存在内部错误，返回 500 状态码
	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
