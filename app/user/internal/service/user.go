package service

import (
	"context"
	"sync"

	dao "github.com/leexsh/go-todolist/app/user/internal/repository"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/myerr"
)

var (
	UserSvrInstance *UserService
	UserSvrOnce     *sync.Once
)

type UserService struct {
	pb_todolist.UnimplementedTaskServiceServer
}

func GetUserService() *UserService {
	UserSvrOnce.Do(func() {
		UserSvrInstance = &UserService{}
	})
	return UserSvrInstance
}

func (s *UserService) UserLogin(ctx context.Context, req *pb_todolist.UserRequest) (*pb_todolist.UserDetailResponse, error) {
	resp := new(pb_todolist.UserDetailResponse)
	resp.Code = myerr.SUCCESS
	r, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = myerr.ERROR
		return resp, err
	}
	resp.UserDetail = &pb_todolist.UserResponse{
		UserID:   r.UserID,
		UserName: r.UserName,
		NickName: r.NickName,
	}
	return resp, nil
}

func (s *UserService) UserRegister(ctx context.Context, req *pb_todolist.UserRequest) (*pb_todolist.UserCommonResposne, error) {
	resp := new(pb_todolist.UserCommonResposne)
	resp.Code = myerr.SUCCESS
	err := dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = myerr.ERROR
		return resp, err
	}
	resp.Meg = myerr.GetMsg(myerr.SUCCESS)
	return resp, nil
}
func (s *UserService) UserLogout(ctx context.Context, req *pb_todolist.UserRequest) (*pb_todolist.UserCommonResposne, error) {
	resp := new(pb_todolist.UserCommonResposne)

	return resp, nil
}
