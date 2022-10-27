package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

// TimeoutHandler 超时的中间件
func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	// 使用函数回调
	return func(c *Context) error {
		finish := make(chan struct{}, 1)       // 这个 channel 负责通知结束
		panicChan := make(chan interface{}, 1) // 这个 channel 负责通知 panic 异常

		// 执行业务逻辑前预操作：初始化超时 context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 执行具体的业务逻辑
			fun(c) // Timeout 中间件以下一层的 ControllerHandler 作为参数

			finish <- struct{}{}
		}()

		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			log.Println(p)
			c.responseWriter.WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}
