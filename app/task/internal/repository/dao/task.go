package dao

import (
	"context"

	"github.com/leexsh/go-todolist/app/task/internal/repository/model"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"gorm.io/gorm"
)

type TaskDao struct {
	*gorm.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	return &TaskDao{NewDBClient(ctx)}
}

func (t *TaskDao) CreateTask(req *pb_todolist.TaskRequest) error {
	task := &model.Task{
		TaskId:    req.TaskID,
		UserId:    req.UserID,
		Status:    req.Status,
		Title:     req.Title,
		Content:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	err := t.Model(&model.Task{}).Create(task)
	if err != nil {
		return err.Error
	}
	return nil
}

func (t *TaskDao) ListTaskByUserId(userId int64) (r []*model.Task, err error) {
	err = t.Model(&model.Task{}).Where("user_id=?", userId).Find(&r).Error
	return
}

func (t *TaskDao) DeleteTaskById(taskId, userId int64) error {
	err := t.Model(&model.Task{}).Where("task_id=? AND user_id=?", taskId, userId).Delete(&model.Task{}).Error
	return err
}
func (t *TaskDao) UpdateTask(req *pb_todolist.TaskRequest) (err error) {
	taskUpdateMap := make(map[string]interface{})
	taskUpdateMap["title"] = req.Title
	taskUpdateMap["content"] = req.Content
	taskUpdateMap["status"] = int(req.Status)
	taskUpdateMap["start_time"] = req.StartTime
	taskUpdateMap["end_time"] = req.EndTime
	err = t.Model(&model.Task{}).
		Where("task_id=?", req.TaskID).Updates(&taskUpdateMap).Error

	return
}
