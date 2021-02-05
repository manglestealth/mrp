package main

import (
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/models"
)

var ProxyServers map[string]*models.ProxyServer = make(map[string]*models.ProxyServer)
func main(){
	commonConfig := &config.Config{}
	commonConfig, ProxyServers = config.LoadConf("frpc")
}
