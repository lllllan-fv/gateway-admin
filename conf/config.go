package conf

import (
	"net"
	"time"
)

type Config struct {
	MySQL
	Redis
}

type MySQL struct {
	Addr     string `mapstructure:"addr"     default:"localhost:3306"`
	User     string `mapstructure:"user"     default:"root"`
	DBName   string `mapstructure:"db_name"  default:"goadmin"`
	Password string `mapstructure:"password" default:"123456"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"     default:"localhost:6379"`
	Password string `mapstructure:"password" default:""`
	Prefix   string `mapstructure:""         default:""`
	DB       uint   `mapstructure:"db"       default:"0"`
}

var err error
var config Config
var TimeLocation *time.Location
var TimeFormat = "2006-01-02 15:04:05"
var DateFormat = "2006-01-02"
var LocalIP = net.ParseIP("127.0.0.1")

func GetConfig() *Config {
	return &config
}

func init() {
	config = Config{
		MySQL: MySQL{
			Addr:     "localhost:3306",
			User:     "root",
			DBName:   "goadmin",
			Password: "123456",
		},
		Redis: Redis{
			Addr: "localhost:6379",
		},
	}

	// 设置时区
	if TimeLocation, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
}
