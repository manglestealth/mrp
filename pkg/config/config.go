package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	BindAddr	string
	BindPort	int64
	LogFile		string
	LogLevel	string
	LogWay		string
}

func LoadConf(path string) *Config {
	viper.SetConfigType("toml")
	viper.AddConfigPath("../../conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	err := viper.ReadInConfig()
	if err != nil{
		log.Fatal("read config failed %v\n", err)
	}
	return &Config{
		BindAddr: viper.GetString("common.bind_addr"),
		BindPort: viper.GetInt64("common.bind_port"),
		LogFile:  viper.GetString("common.log_file"),
		LogLevel: viper.GetString("common.log_level"),
		LogWay:   viper.GetString("common.log_way"),
	}

}
