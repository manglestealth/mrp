package main

import (
	"encoding/json"
	"github.com/manglestealth/mrp/pkg/conn"
	"github.com/manglestealth/mrp/pkg/models"
	"log"
)

var ProxyServers map[string]*models.ProxyServer = make(map[string]*models.ProxyServer)

func ProcessControlConn(l *conn.Listener){
	for{
		c := l.GetConn()
		log.Printf("Get one new conn %v\n", c)
		go controlWorker(c)
	}
}


func controlWorker(c *conn.Conn){
	//读取客户端发送给服务器的第一条信息，失败则关闭连接
	res, err := c.ReadLine()
	if err != nil{
		log.Fatalf("Read error %v\n", err)
	}
	log.Printf("get %s", res)

	clientCtlReq := &models.ClientCtlReq{}
	clientCtlRes := &models.ClientCtlRes{}

	if err := json.Unmarshal([]byte(res), &clientCtlReq); err != nil {
		log.Fatalf("Parse err : %v : %s", err, res)
	}

	succ, msg, needRes := checkProxy(clientCtlReq, c)
	if !succ{
		clientCtlRes.Code = 1
		clientCtlRes.Msg = msg
	}
}

func checkProxy(req *models.ClientCtlReq, c *conn.Conn)(bool, string, bool){
	succ := false
	needRes := true
}