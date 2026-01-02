package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	// 创建一个默认的 Hertz 服务器实例（默认使用 8888 端口）
	// 使用 server.Default() 会默认添加一些中间件，例如恢复中间件
	// 使用 server.New() 则可以创建一个没有任何默认中间件的纯净实例
	h := server.Default()

	// 注册一个 GET 路由处理函数
	// 第一个参数是路由路径，第二个参数是处理函数
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		// 返回 JSON 响应
		ctx.JSON(consts.StatusOK, map[string]interface{}{
			"message": "pong",
			"status":  "success",
		})
	})

	// 启动服务器
	// Spin 方法会阻塞当前 goroutine，直到服务器被关闭
	h.Spin()
}
