package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Context 自定义 Context 结构
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context

	// 是否超时标记位
	hasTimeout bool
	// 写保护机制
	writeMux *sync.Mutex

	// 当前请求的 handler 链条
	handlers []ControllerHandler
	index    int // 当前请求调用到调用链的哪个节点

	params map[string]string // url 路由匹配的参数
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	// index 初始值应该为 -1，每次调用都会自增1，这样才能保证第一次调用的时候 index 为 0
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writeMux:       &sync.Mutex{},
		index:          -1,
	}
}

// #region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writeMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// SetHandlers 为 context 设置 handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

// SetParams 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// Next 核心函数，调用 context 的下一个函数
/*
	Next() 方法每调用一次，就将这个控制器链路的调用控制器，往后移动一位
	Next() 函数会在框架的两个地方被调用：
		1.在此请求处理的入口处，即 core 的 ServeHttp
		2.在每个中间件的逻辑代码中，用于调用下个中间件
*/
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		} // func(c *Context) error
	}
	return nil
}

// #endregion

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// #region implement context.Context

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// #endregion
