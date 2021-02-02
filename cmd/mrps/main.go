package main

import (
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/conn"
	"github.com/manglestealth/mrp/pkg/models"
	"log"
)

var ProxyServers map[string]*models.ProxyServer = make(map[string]*models.ProxyServer)

func main() {
	//frpConf := config.LoadConf("frps")
	commonConfig, proxyServers := config.LoadConf("frps")
	//fmt.Println(frpConf)
	l, err := conn.Listen(commonConfig.BindAddr, commonConfig.BindPort)
	if err != nil {
		log.Fatalf("create listener error %v", err)
	}
	c := l.GetConn()
	controlWorker(c)
}
