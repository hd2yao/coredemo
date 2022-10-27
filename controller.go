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
	}()
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

	return nil
}
