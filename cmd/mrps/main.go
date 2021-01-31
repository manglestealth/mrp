package main

import (
	"fmt"
	"github.com/manglestealth/mrp/pkg/conn"
	"log"
)

func main() {
	//frpConf := config.LoadConf("frps")
	//fmt.Println(frpConf)
	l,err  := conn.Listen("0.0.0.0", 3000)
	if err != nil{
		log.Fatal(err)
	}
	c := l.GetConn()
	fmt.Println(c.GetLocalAddr())
}
