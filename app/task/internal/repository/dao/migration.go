package dao

import (
	"os"

	"github.com/leexsh/go-todolist/app/task/internal/repository/model"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Task{},
		)
	if err != nil {
		os.Exit(0)
	}
}
