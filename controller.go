package main

import (
	"context"
	"fmt"
	"goweb/framework"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context) error {

	// 1.生成一个超时的 Context
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	// 2.创建一个新的 Goroutine 来处理业务逻辑
	finish := make(chan struct{}, 1)       // 这个 channel 负责通知结束
	panicChan := make(chan interface{}, 1) // 这个 channel 负责通知 panic 异常

	go func() {
		// 异常处理
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		// 这里做具体的业务
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")

		// 新的 goroutine 结束时通过一个 finish 通道告知 父 goroutine
		finish <- struct{}{}

		// 3.监听
		select {
		// 监听 panic
		case p := <-panicChan:
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			log.Println(p)
			c.Json(500, "panic")
		// 监听结束事件
		case <-finish:
			fmt.Println("finish")
		// 监听超时事件
		case <-durationCtx.Done():
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			c.Json(500, "time out")
			c.SetHasTimeout()
		}
	}()

	return nil
}

//// Foo1 控制器(未封装)
//func Foo1(request *http.Request, response http.ResponseWriter) {
//	obj := map[string]interface{}{
//		"data": nil,
//	}
//
//	// 设置控制器的 response 的 header 部分
//	response.Header().Set("Content-Type", "application/json")
//
//	// 从请求体中获取参数
//	foo := request.PostFormValue("foo")
//	if foo == "" {
//		foo = "10"
//	}
//	fooInt, err := strconv.Atoi(foo)
//	if err != nil {
//		response.WriteHeader(500)
//		return
//	}
//
//	// 构建返回结构
//	obj["data"] = fooInt
//	byt, err := json.Marshal(obj)
//	if err != nil {
//		response.WriteHeader(500)
//		return
//	}
//
//	// 构建返回状态，输出返回结构
//	response.WriteHeader(200)
//	response.Write(byt)
//	return
//}
//
//// Foo2 控制器(封装)
//func Foo2(ctx *framework.Context) error {
//	obj := map[string]interface{}{
//		"data": nil,
//	}
//
//	// 从请求体中获取参数
//	fooInt := ctx.FormInt("foo", 10)
//	// 构建返回结构
//	obj["data"] = fooInt
//	// 输出返回结构
//	return ctx.Json(http.StatusOK, obj)
//}
//
//func Foo3(ctx *framework.Context) error {
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//	return rdb.Set(ctx, "key", "value", 0).Err()
//}
