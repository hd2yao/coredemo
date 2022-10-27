package middleware

import "goweb/framework"

// Recovery recovery 机制，将协程中的函数异常进行捕获
func Recovery() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		// 核心在增加这个 recover 机制，捕获 c.Next() 出现的 panic
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		// 使用 next() 执行具体的业务逻辑
		c.Next()
		return nil
	}
}
