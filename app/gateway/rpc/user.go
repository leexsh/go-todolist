package rpc

import (
	"context"
	"errors"

	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/myerr"
)

func UserLogin(ctx context.Context, req *pb_todolist.UserRequest) (resp *pb_todolist.UserDetailResponse, err error) {
	resp, err = UserClient.UserLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New("err")
	}
	return
}

func UserRegister(ctx context.Context, req *pb_todolist.UserRequest) (resp *pb_todolist.UserCommonResposne, err error) {
	resp, err = UserClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New(resp.Meg)
	}
	return
}
