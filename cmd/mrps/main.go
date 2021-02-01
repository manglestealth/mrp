package main

import (
	"fmt"
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/conn"
	"log"
)

func main() {
	frpConf := config.LoadConf("frps")
	fmt.Println(frpConf)
	l,err  := conn.Listen(frpConf.BindAddr, frpConf.BindPort)

}
