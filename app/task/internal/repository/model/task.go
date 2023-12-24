package model

type Task struct {
	TaskId    int64 `gorm:"primarykey"`
	UserId    int64 `gorm:"index"`
	Status    int64 `gorm:"default:0"`
	Title     string
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
}

func (t *Task) Table() string {
	return "task"
}
