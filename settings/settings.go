package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 全局配置
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MySqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
}

type MySqlConfig struct {
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbname"`
	Charset     string `mapstructure:"charset"`
	Port        int    `mapstructure:"port"`
	MaxConn     int    `mapstructure:"max_conns"`
	MaxIdleConn int    `mapstructure:"max_idle_conns"`
}

func Init(path string) (err error) {
	viper.SetConfigFile(path)
	//viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项, 一般用于拉取配置中心配置类型
	viper.AddConfigPath(".")   // 还可以在工作目录中查找配置
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误

		fmt.Printf("viper.ReadInConfig() failed, error: %v\n", err)
		return
	}

	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err: %v", err)
	}

	viper.WatchConfig() // 热加载配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config update...")

		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err: %v", err)
		}
	})
	return
}
