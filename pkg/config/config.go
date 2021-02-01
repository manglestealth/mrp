package config

import (
	"fmt"
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

func LoadConf(path string){
	viper.SetConfigType("toml")
	viper.AddConfigPath("../../conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	err := viper.ReadInConfig()
	if err != nil{
		log.Fatal("read config failed %v\n", err)
	}
	commonConfig := &Config{}
	for k,j := range viper.AllSettings(){
		fmt.Println(k,j)
		if k == "common"{
			commonConfig.BindAddr = j["bind_addr"]
			commonConfig.BindPort= viper.GetInt64("common.bind_port")
			commonConfig.LogFile= viper.GetString("common.log_file")
			commonConfig.LogLevel= viper.GetString("common.log_level")
			commonConfig.LogWay=   viper.GetString("common.log_way")


		}else{
				fmt.Println(2)
		}
	}
}
