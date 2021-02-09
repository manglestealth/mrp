package main

import (
	"encoding/json"
	"github.com/manglestealth/mrp/pkg/config"
	"github.com/manglestealth/mrp/pkg/conn"
	"github.com/manglestealth/mrp/pkg/models"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
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

	c := &conn.Conn{}
	err := c.ConnServer(commonConfig.BindAddr, commonConfig.BindPort)
	if err != nil{
		log.Fatalf("Proxyname [%s], connect to server [%s:%d] error, %v", client.Name, commonConfig.BindAddr, commonConfig.BindPort, err)
	}

	defer c.Close()

	req := &models.ClientCtlReq{
		Type: models.ControlConn,
		ProxyName: client.Name,
		Passwd: client.Passwd,
	}

	buf, _ := json.Marshal(req)
	err = c.Write(string(buf) + "\n")
	if err != nil{
		log.Fatalf("ProxyName [%s], write to server error, %v", client.Name, err)
		return
	}

	res, err := c.ReadLine()
	if err != nil{
		log.Fatalf("ProxyName [%s], read from server error, %v", client.Name, err)
	}

	log.Infof("ProxyName [%s], read [%s]", client.Name, res)

	clientCtlRes := &models.ClientCtlRes{}
	if err = json.Unmarshal([]byte(res), &clientCtlRes); err != nil {
		log.Fatalf("ProxyName [%s], format server response error, %v", client.Name, err)
		return
	}

	if clientCtlRes.Code != 0 {
		log.Fatalf("ProxyName [%s], start proxy error, %s", client.Name, clientCtlRes.Msg)
		return
	}

	for {
		// ignore response content now
		_, err := c.ReadLine()
		if err == io.EOF {
			log.Infof("ProxyName [%s], server close this control conn", client.Name)
			break
		} else if err != nil {
			log.Warnf("ProxyName [%s], read from server error, %v", client.Name, err)
			continue
		}
		client.StartTunnel(commonConfig.BindAddr, commonConfig.BindPort)
	}

}
