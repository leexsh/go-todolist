package main

import (
	"github.com/leexsh/go-todolist/app/task/internal/repository/dao"
	"github.com/leexsh/go-todolist/config"
)

func main() {
	config.InitConfig()
	dao.InitDB()
	// etcdAdd := []string{config.Conf.Etcd.Address}

}
