package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leexsh/go-todolist/app/gateway/rpc"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/ctl"
)

func GetTaskList(ctx *gin.Context) {
	var req pb_todolist.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "err"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func CreateTask(ctx *gin.Context) {
	var req pb_todolist.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskCreate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "err"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func UpdateTask(ctx *gin.Context) {
	var req pb_todolist.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskUpdate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "err"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func DeleteTask(ctx *gin.Context) {
	var req pb_todolist.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskDelete(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "err"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
