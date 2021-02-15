package main

import (
	"encoding/json"
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/conn"
	"github.com/manglestealth/mrp/pkg/models"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

const(
	//重连最小间隔时间
	sleepMinDuration = 1
	sleepMaxDuration = 60
)

var ProxyClients map[string]*models.ProxyClient = make(map[string]*models.ProxyClient)
var commonConfig = &config.Config{}
func main(){
	//
	commonConfig, ProxyClients = config.LoadClientConf("frpc")

	var wait sync.WaitGroup
	wait.Add(len(ProxyClients))
	for _, client := range ProxyClients{
		go controlProcess(client, &wait)
	}
	wait.Wait()
}

func controlProcess(client *models.ProxyClient, wait *sync.WaitGroup){
	defer wait.Done()

	c := loginToServer(client)
	if c == nil{
		log.Fatalf("Proxyname [%s], connect to server failed", client.Name)
		return
	}
	defer c.Close()

	for{
		_,err := c.ReadLine()
		if err == io.EOF{
			log.Debugf("ProxyName [%s], server close this control conn", client.Name)
			var sleepTime time.Duration = 1
			for{
				log.Debugf("ProxyName [%s], try to reconnect to server[%s:%d]...", client.Name, commonConfig.BindAddr, commonConfig.BindPort)
				tmpConn := loginToServer(client)
				if tmpConn != nil{
					c.Close()
					c = tmpConn
					break
				}

				if sleepTime < 60{
					sleepTime++
				}
				time.Sleep(sleepTime * time.Second)
			}
			continue
		}else if err != nil{
			log.Warnf("ProxyName [%s], read from server error, %v", client.Name, err)
			continue
		}

		client.StartTunnel(commonConfig.BindAddr, commonConfig.BindPort)
	}
}

func loginToServer(cli *models.ProxyClient) (connection *conn.Conn){
	c := &conn.Conn{}

	connection = nil
	for i := 0; i < 1; i++{
		 err := c.ConnServer(commonConfig.BindAddr, commonConfig.BindPort)
		 if err != nil{
		 	log.Warnf("ProxyName [%s], connect to server [%s:%d] error, %v", cli.Name, commonConfig.BindAddr, commonConfig.BindPort, err)
		 	break
		 }

		 req := &models.ClientCtlReq{
		 	Type: models.ControlConn,
		 	ProxyName: cli.Name,
		 	Passwd: cli.Passwd,
		 }

		 buf, _ := json.Marshal(req)
		 err = c.Write(string(buf) + "\n")
		 if err != nil{
		 	log.Warnf("ProxyName [%s], write to server error, %v", cli.Name, err)
		 	break
		 }

		 res, err := c.ReadLine()
		 if err != nil{
		 	log.Warnf("ProxyName [%s], read from server error, %v", cli.Name, err)
		 	break
		 }

		 log.Debugf("ProxyName [%s], read [%s]", cli.Name, res)
		 clientCtlRes := &models.ClientCtlRes{}
		 if err = json.Unmarshal([]byte(res), &clientCtlRes); err != nil{
		 	log.Warnf("ProxyName [%s], format server response error, %v", cli.Name, err)
		 	break
		 }

		 if clientCtlRes.Code != 0{
		 	log.Warnf("ProxyName [%s], start proxy error, %s", cli.Name, clientCtlRes.Msg)
		 	break
		 }

		 connection = c
	}

	if connection == nil{
		c.Close()
	}

	return
}
