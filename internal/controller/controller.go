// Package controller 提供API接口的具体逻辑信息

package controller

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/thunes/pkg/restapi"
)

var (
	_idGenerator, _ = snowflake.NewNode(1)
)

func ginSuccess(ctx *gin.Context, msg string, data ...interface{}) {
	if len(msg) == 0 {
		msg = "success"
	}

	var retData interface{} = data
	if len(data) == 1 {
		retData = data[0]
	}

	ctx.JSON(http.StatusOK, restapi.Response{
		Code: http.StatusOK,
		Msg:  msg,
		Data: retData,
	})
}

func ginInternalErr(ctx *gin.Context, extra ...interface{}) {
	var (
		msg = "server error"
	)

	for _, info := range extra {
		msg = fmt.Sprintf("%s, %v", msg, info)
	}

	ctx.JSON(http.StatusOK, restapi.Response{
		Code: ErrInternal.Code,
		Msg:  msg,
	})
}

func ginInvalidParam(ctx *gin.Context, extra ...interface{}) {
	var (
		msg = "invalid param"
	)

	for _, info := range extra {
		msg = fmt.Sprintf("%s, %v", msg, info)
	}

	ctx.JSON(http.StatusOK, restapi.Response{
		Code: ErrInvalidParam.Code,
		Msg:  msg,
	})
}

func ginError(ctx *gin.Context, err error) {
	var (
		msg        = "error"
		code int64 = -1
	)

	e, ok := err.(*Error)
	if ok {
		code = e.Code
		msg = e.Error()
	}

	ctx.JSON(http.StatusOK, restapi.Response{
		Code: code,
		Msg:  msg,
	})
}

func ginResponse(ctx *gin.Context, resp restapi.Response) {
	ctx.JSON(http.StatusOK, resp)
}
