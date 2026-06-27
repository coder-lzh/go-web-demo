package handler

import (
	"context"
	"go-web-demo/internal/logger"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"go-web-demo/internal/common/response"
	"go-web-demo/internal/dto"
	"go-web-demo/internal/service"
)

type UserHandler struct{}

var userHandler *UserHandler

func init() {
	userHandler = &UserHandler{}
}

func GetUserHandler() *UserHandler {
	return userHandler
}

func (h *UserHandler) Create(c context.Context, ctx *app.RequestContext) {
	var req dto.UserCreateRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid request: "+err.Error()))
		return
	}

	user, err := service.GetUserService().Create(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, response.Fail(consts.StatusInternalServerError, "create failed: "+err.Error()))
		return
	}

	ctx.JSON(consts.StatusOK, response.Success(user))
}

func (h *UserHandler) Update(c context.Context, ctx *app.RequestContext) {
	var req dto.UserUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid request: "+err.Error()))
		return
	}

	user, err := service.GetUserService().Update(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, response.Fail(consts.StatusInternalServerError, "update failed: "+err.Error()))
		return
	}

	ctx.JSON(consts.StatusOK, response.Success(user))
}

func (h *UserHandler) Delete(c context.Context, ctx *app.RequestContext) {
	var req dto.UserDeleteRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid request: "+err.Error()))
		return
	}

	if err := service.GetUserService().Delete(c, req.ID); err != nil {
		ctx.JSON(consts.StatusInternalServerError, response.Fail(consts.StatusInternalServerError, "delete failed: "+err.Error()))
		return
	}

	ctx.JSON(consts.StatusOK, response.Success(nil))
}

func (h *UserHandler) GetByID(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid id"))
		return
	}

	user, err := service.GetUserService().GetByID(c, uint(id))
	if err != nil {
		ctx.JSON(consts.StatusNotFound, response.Fail(consts.StatusNotFound, err.Error()))
		return
	}

	logger.Info("GetByID 查询用户 userId:%s", idStr)

	ctx.JSON(consts.StatusOK, response.Success(user))
}

func (h *UserHandler) GetList(c context.Context, ctx *app.RequestContext) {
	var req dto.UserPageQueryRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid request: "+err.Error()))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	users, err := service.GetUserService().GetList(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, response.Fail(consts.StatusInternalServerError, "query failed: "+err.Error()))
		return
	}

	ctx.JSON(consts.StatusOK, response.Success(users))
}

func (h *UserHandler) GetPageList(c context.Context, ctx *app.RequestContext) {
	var req dto.UserPageQueryRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, response.Fail(consts.StatusBadRequest, "invalid request: "+err.Error()))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	pageResult, err := service.GetUserService().GetPageList(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, response.Fail(consts.StatusInternalServerError, "query failed: "+err.Error()))
		return
	}

	ctx.JSON(consts.StatusOK, response.Success(pageResult))
}
