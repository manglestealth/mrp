package models

import (
	"encoding/json"
	"github.com/manglestealth/mrp/pkg/conn"
	"log"
)

type ProxyClient struct {
	Name string
	Passwd string
	LocalPort int64
}

func (p *ProxyClient) GetLocalConn()(c *conn.Conn, err error)  {
	c = &conn.Conn{}
	err = c.ConnServer("127.0.0.1", p.LocalPort)
	if err != nil{
		log.Fatalf("ProxyName [%s], conn to local port error, %v", p.Name, err)
	}
	return
}

func (p *ProxyClient) GetRemoteConn(addr string, port int64)(c *conn.Conn, err error){
	c = &conn.Conn{}
	defer func(){
		if err != nil{
			c.Close()
		}
	}()

	err = c.ConnServer(addr, port)
	if err != nil{
		log.Fatalf("ProxyName [%s], connect to server [%s:%d] error, %v", p.Name, addr, port, err)
	}

	req := &ClientCtlReq{
		Type : WorkConn,
		ProxyName : p.Name,
		Passwd : p.Passwd,
	}

	buf, _ := json.Marshal(req)
	err = c.Write(string(buf) + "\n")
	if err != nil{
		log.Fatalf("ProxyName [%s], write to server error, %v", p.Name, err)
	}

	return
}

func(p *ProxyClient) StartTunnel(serverAddr string, serverPort int64)(err error){
	localConn, err := p.GetLocalConn()
	if err != nil{
		return
	}

	remoteConn, err :=  p.GetRemoteConn(serverAddr, serverPort)
	if err != nil{
		return
	}
	log.Printf("Join two conns, (l[%s] r[%s]) (l[%s] r[%s])", localConn.GetLocalAddr(), localConn.GetRemoteAddr(), remoteConn.GetLocalAddr(), remoteConn.GetRemoteAddr())
	go conn.Join(localConn, remoteConn)
    return nil
}