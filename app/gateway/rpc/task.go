package rpc

import (
	"context"
	"errors"

	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/myerr"
)

func TaskCreate(ctx context.Context, req *pb_todolist.TaskRequest) (resp *pb_todolist.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New(resp.Msg)
	}
	return
}

func TaskUpdate(ctx context.Context, req *pb_todolist.TaskRequest) (resp *pb_todolist.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New(resp.Msg)
	}
	return
}

func TaskDelete(ctx context.Context, req *pb_todolist.TaskRequest) (resp *pb_todolist.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskDelete(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New(resp.Msg)
	}
	return
}

func TaskList(ctx context.Context, req *pb_todolist.TaskRequest) (resp *pb_todolist.TasksDetailResponse, err error) {
	resp, err = TaskClient.TaskShow(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != myerr.SUCCESS {
		err = errors.New("error")
	}
	return
}
