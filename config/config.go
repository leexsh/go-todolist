package config

import (
	"os"

	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Server   *Server             `yaml:"server"`
	Mysql    *Mysql              `yaml:"mysql"`
	Redis    *Redis              `yaml:"redis"`
	Etcd     *Etcd               `yaml:"etcd"`
	Services map[string]*Service `yaml:"services"`
	Domains  map[string]*Domain  `yaml:"domain"`
}

type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"version"`
	JwtSecret string `yaml:"jwtSecret"`
}

type Mysql struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DataBase   string `yaml:"dataBase"`
	UserName   string `yaml:"userName"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
}

type Redis struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Domain struct {
	Name string `yaml:"name"`
}

type Service struct {
	Name        string   `yaml:"name"`
	LoadBalance bool     `yaml:"loadBalance"`
	Addr        []string `yaml:"addr"`
}

func InitConfig() {
	viperConfig := viper.New()
	viperConfig.SetConfigName("config")
	viperConfig.SetConfigType("yaml")
	workdir, _ := os.Getwd()
	viperConfig.AddConfigPath(workdir + "/config")
	err := viperConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viperConfig.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}

}
