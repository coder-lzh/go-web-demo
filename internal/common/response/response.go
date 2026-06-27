package response

import (
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageResult[T any] struct {
	List       []T   `json:"list"`
	Total      int64 `json:"total"`
	Page       int32 `json:"page"`
	PageSize   int32 `json:"page_size"`
	TotalPages int32 `json:"total_pages"`
}

func Success(data interface{}) *Result {
	return &Result{
		Code:    consts.StatusOK,
		Message: "success",
		Data:    data,
	}
}

func Fail(code int, message string) *Result {
	return &Result{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func FailWithData(code int, message string, data interface{}) *Result {
	return &Result{
		Code:    code,
		Message: message,
		Data:    data,
	}
}