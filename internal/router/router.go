package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"go-web-demo/internal/config"
	"go-web-demo/internal/handler"
	"go-web-demo/internal/logger"
	"go-web-demo/internal/middleware"
)

func SetupRoutes(h *server.Hertz, cfg *config.Config) {
	h.Use(middleware.Logging())

	h.GET("/ping", pingHandler(cfg))
	h.GET("/config", configHandler(cfg))

	userGroup := h.Group("/api/v1/users")
	{
		userGroup.POST("/create", handler.GetUserHandler().Create)
		userGroup.POST("/update", handler.GetUserHandler().Update)
		userGroup.POST("/delete", handler.GetUserHandler().Delete)
		userGroup.GET("/detail/:id", handler.GetUserHandler().GetByID)
		userGroup.GET("/page", handler.GetUserHandler().GetPageList)
		userGroup.GET("/list", handler.GetUserHandler().GetList)
	}

	logger.Infof("Server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
}

func pingHandler(cfg *config.Config) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, map[string]interface{}{
			"message": "pong",
			"status":  "success",
			"env":     cfg.Env,
		})
	}
}

func configHandler(cfg *config.Config) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, map[string]interface{}{
			"app":    cfg.App,
			"server": cfg.Server,
			"env":    cfg.Env,
		})
	}
}
