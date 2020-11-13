// Package settings provides ...
package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存应用程序所有的配置信息
var Conf = new(AppConfig)

// AppConfig 应用程序配置信息
type AppConfig struct {
	Name      string `mapstructure:",name"`
	Mode      string `mapstructure:",mode"`
	Version   string `mapstructure:",version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

// LogConfig zap配置信息
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backup"`
}

// MySQLConfig MySQL配置信息
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idel_conns"`
}

// RedisConfig Redis配置信息
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Init 初始化项目配置
func Init() (err error) {
	// 指定配置文件名称
	viper.SetConfigName("config")
	// 指定配置文件家在路径
	viper.AddConfigPath("./conf")
	// 读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 将读取的配置信息保存到全局变量 Conf 中
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("viper.Unmarshal conf failed, err:%v\n", err))
	}

	// 配置热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件更新了 ...\n")
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("viper.Unmarshal conf failed, err:%v\n", err))
		}
		fmt.Printf("Conf 成功同步程序最新配置 ...\n")
	})

	return
}
