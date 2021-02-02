package config

import (
	"github.com/manglestealth/mrp/pkg/models"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	BindAddr string
	BindPort int64
	LogFile  string
	LogLevel string
	LogWay   string
}

var ProxyServers map[string]*models.ProxyServer = make(map[string]*models.ProxyServer)

func LoadConf(path string) (*Config, map[string]*models.ProxyServer) {
	viper.SetConfigType("toml")
	viper.AddConfigPath("../../conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed %v\n", err)
	}

	commonConfig := &Config{}
	for k, j := range viper.AllSettings() {
		//fmt.Println(k,j)
		m := j.(map[string]interface{})
		if k == "common" {

			commonConfig.BindAddr = m["bind_addr"].(string)
			commonConfig.BindPort = m["bind_port"].(int64)
			commonConfig.LogFile = m["log_file"].(string)
			commonConfig.LogLevel = m["log_level"].(string)
			commonConfig.LogWay = m["log_way"].(string)
			//fmt.Println(commonConfig)
		} else {

			for name, section := range m {
				sectionMap := section.(map[string]interface{})
				proxyServer := &models.ProxyServer{}
				proxyServer.Name = name
				proxyServer.Passwd = sectionMap["passwd"].(string)
				proxyServer.Init()
				ProxyServers[name] = proxyServer
			}
		}
	}
	return commonConfig, ProxyServers
}
