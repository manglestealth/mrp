package main

import (
	"encoding/json"
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/conn"
	"github.com/manglestealth/mrp/pkg/models"
	"log"
	"fmt"
)

var ProxyServers map[string]*models.ProxyServer = make(map[string]*models.ProxyServer)

func main() {
	//frpConf := config.LoadConf("frps")
	commonConfig := &config.Config{}
	commonConfig, ProxyServers = config.LoadConf("frps")
	//fmt.Println(frpConf)
	l, err := conn.Listen(commonConfig.BindAddr, commonConfig.BindPort)
	if err != nil {
		log.Fatalf("create listener error %v", err)
	}
	c := l.GetConn()
	controlWorker(c)
}

func ProcessControlConn(l *conn.Listener) {
	for {
		c := l.GetConn()
		log.Printf("Get one new conn %v\n", c)
		go controlWorker(c)
	}
}

func controlWorker(c *conn.Conn) {
	//读取客户端发送给服务器的第一条信息，失败则关闭连接
	res, err := c.ReadLine()
	if err != nil {
		log.Fatalf("Read error %v\n", err)
	}
	log.Printf("get %s", res)

	clientCtlReq := &models.ClientCtlReq{}
	clientCtlRes := &models.ClientCtlRes{}

	if err := json.Unmarshal([]byte(res), &clientCtlReq); err != nil {
		log.Fatalf("Parse err : %v : %s", err, res)
	}

	succ, msg, needRes := checkProxy(clientCtlReq, c)
	if !succ {
		clientCtlRes.Code = 1
		clientCtlRes.Msg = msg
	}

	if needRes{
		buf, _ := json.Marshal(clientCtlRes)
		err = c.Write(string(buf) + "\n")
		if err != nil{
			log.Fatalf("write error, %v", err)
		}
	}else{
		return
	}

	defer c.Close()

	server, ok := ProxyServers[clientCtlReq.ProxyName]
	if !ok{
		log.Fatalf("ProxyName [%s] is not exist", clientCtlReq.ProxyName)
	}
	serverCtlReq := &models.ClientCtlReq{}
	serverCtlReq.Type = models.WorkConn

	for {
		server.WaitUserConn()
		buf, _ := json.Marshal(serverCtlReq)
		err = c.Write(string(buf) + "\n")
		if err != nil{
			log.Fatalf("ProxyName [%s], write to client error, proxy exit", server.Name)
			server.Close()
			return
		}
	}

	return
}

func checkProxy(req *models.ClientCtlReq, c *conn.Conn) (succ bool, msg string, needRes bool) {
	succ = false
	needRes = true

	server, ok := ProxyServers[req.ProxyName]
	if !ok {
		msg = fmt.Sprintf("ProxyName [%s] is not exist", req.ProxyName)
		log.Fatal(msg)
		return
	}

	if req.Passwd != server.Passwd {
		msg = fmt.Sprintf("ProxyName [%s], password is not correct", req.ProxyName)
		log.Fatal(msg)
		return
	}

	if req.Type == models.ControlConn {
		if server.Status != models.Idle {
			msg = fmt.Sprintf("ProxyName [%s], already in use", req.ProxyName)
			log.Fatal(msg)
			return
		}

		err := server.Start()
		if err != nil {
			msg = fmt.Sprintf("ProxyName [%s], start proxy error: %v", req.ProxyName, err.Error())
			log.Fatal(msg)
			return
		}
	} else if req.Type == models.WorkConn {
		needRes = false
		if server.Status != models.Working {
			msg = fmt.Sprintf("ProxyName [%s], is not working when it gets one new work conn", req.ProxyName)
			log.Print(msg)
			return
		}

		server.CliConnChan <- c
	} else {
		msg = fmt.Sprintf("ProxyName [%s], type unsupport", req.ProxyName)
		log.Print(msg)
		return
	}

	succ = true
	return
}
