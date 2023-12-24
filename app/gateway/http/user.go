package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leexsh/go-todolist/app/gateway/rpc"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/ctl"
	"github.com/leexsh/go-todolist/pkg/res"
	"github.com/leexsh/go-todolist/util/jwt"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq pb_todolist.UserRequest
	if err := ctx.Bind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	r, err := rpc.UserRegister(ctx, &userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var req pb_todolist.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC服务调用错误"))
		return
	}

	token, err := jwt.GenerateToken(123)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "加密错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res.TokenData{User: userResp, Token: token}))
}
