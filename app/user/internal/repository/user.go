package dao

import (
	"context"
	"errors"

	"github.com/leexsh/go-todolist/app/user/internal/model"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *pb_todolist.UserRequest) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", req.UserName).
		First(&r).Error

	return
}

// CreateUser 用户创建
func (dao *UserDao) CreateUser(req *pb_todolist.UserRequest) (err error) {
	var user model.User
	var count int64
	dao.Model(&model.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count != 0 {
		return errors.New("UserName Exist")
	}

	user = model.User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	_ = user.SetPassword(req.Passwd)
	if err = dao.Model(&model.User{}).Create(&user).Error; err != nil {
		return
	}

	return
}
