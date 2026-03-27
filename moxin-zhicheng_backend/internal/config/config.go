package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppMode  string         `mapstructure:"app_mode"`
	Database DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

var Conf = &Config{}

func InitConfig() {
	// 注意：这里要匹配你截图里的文件名 config_dev
	viper.SetConfigName("config_dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // 确保 configs 文件夹在项目根目录

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}
}
