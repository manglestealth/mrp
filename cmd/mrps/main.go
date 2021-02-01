package main

import (
	"github.com/manglestealth/mrp/pkg/config"
)

func main() {
	//frpConf := config.LoadConf("frps")
	config.LoadConf("frps")
	//fmt.Println(frpConf)
	//l,err  := conn.Listen(frpConf.BindAddr, frpConf.BindPort)

}
