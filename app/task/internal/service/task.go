package service

import (
	"context"
	"sync"

	"github.com/leexsh/go-todolist/app/task/internal/repository/dao"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/myerr"
)

var (
	TaskSvrInstance *TaskService
	TaskSvrOnce     *sync.Once
)

type TaskService struct {
	pb_todolist.UnimplementedTaskServiceServer
}

func GetTaskService() *TaskService {
	TaskSvrOnce.Do(func() {
		TaskSvrInstance = &TaskService{}
	})
	return TaskSvrInstance
}

func (*TaskService) TaskCreate(ctx context.Context, req *pb_todolist.TaskRequest) (*pb_todolist.TaskCommonResponse, error) {
	resp := new(pb_todolist.TaskCommonResponse)
	resp.Code = myerr.SUCCESS
	e := dao.NewTaskDao(ctx).CreateTask(req)
	if e != nil {
		resp.Code = myerr.ERROR
		resp.Msg = myerr.GetMsg(myerr.ERROR)
		resp.Data = e.Error()
		return resp, e
	}
	resp.Msg = myerr.GetMsg(int(resp.Code))
	return resp, nil
}

func (*TaskService) TaskUpdate(ctx context.Context, req *pb_todolist.TaskRequest) (*pb_todolist.TaskCommonResponse, error) {
	resp := new(pb_todolist.TaskCommonResponse)
	resp.Code = myerr.SUCCESS
	err := dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = myerr.ERROR
		resp.Msg = myerr.GetMsg(myerr.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = myerr.GetMsg(myerr.SUCCESS)
	return resp, nil
}

func (*TaskService) TaskShow(ctx context.Context, req *pb_todolist.TaskRequest) (*pb_todolist.TasksDetailResponse, error) {
	resp := new(pb_todolist.TasksDetailResponse)
	rtask, err := dao.NewTaskDao(ctx).ListTaskByUserId(req.UserID)
	if err != nil {
		resp.Code = myerr.ERROR
		return resp, err
	}
	for i := 0; i < len(rtask); i++ {
		resp.TasksDetails = append(resp.TasksDetails, &pb_todolist.TaskModel{
			TaskID:    rtask[i].TaskId,
			UserID:    rtask[i].UserId,
			Status:    rtask[i].Status,
			Title:     rtask[i].Title,
			Content:   rtask[i].Content,
			StartTime: rtask[i].StartTime,
			EndTime:   rtask[i].EndTime,
		})
	}
	resp.Code = myerr.SUCCESS
	return resp, nil
}

func (*TaskService) TaskDelete(ctx context.Context, req *pb_todolist.TaskRequest) (*pb_todolist.TaskCommonResponse, error) {
	resp := new(pb_todolist.TaskCommonResponse)
	resp.Code = myerr.SUCCESS
	err := dao.NewTaskDao(ctx).DeleteTaskById(req.TaskID, req.UserID)
	if err != nil {
		resp.Code = myerr.ERROR
		resp.Msg = myerr.GetMsg(myerr.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = myerr.GetMsg(myerr.ERROR)
	return resp, nil
}
