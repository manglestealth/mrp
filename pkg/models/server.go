package models

import (
	"container/list"
	"github.com/manglestealth/mrp/pkg/conn"
	"log"
	"sync"
)

const (
	Idle = iota
	Working
)

type ProxyServer struct {
	Name       string
	Passwd     string
	BindAddr   string
	ListenPort int64

	Status   int64
	Listener *conn.Listener //接受一个远程用户的连接请求
	CtlMsgChan chan int64 //每次有新的用户连接时，向通道发送一条信息
	CliConnChan chan *conn.Conn //获取用户连接的通道
	UserConnList *list.List //保存用户的连接
	Mutex sync.Mutex
}

func(p *ProxyServer)Init(){
	p.Status = Idle //默认闲置
	p.CtlMsgChan = make(chan int64)
	p.CliConnChan = make(chan *conn.Conn)
	p.UserConnList = list.New()
}

func(p *ProxyServer)Lock(){
	p.Mutex.Lock()
}

func(p *ProxyServer)Unlock(){
	p.Mutex.Unlock()
}

func(p *ProxyServer)Start() (err error){
	p.Listener, err = conn.Listen(p.BindAddr, p.ListenPort)
	if err != nil{
		return err
	}
	p.Status = Working

	//开启监听
	go func(){
		for {
			//阻塞
			c := p.Listener.GetConn()
			log.Printf("ProxyName [%s], get one new user conn [%s]", p.Name, c.GetRemoteAddr())

			//加入队列
			p.Lock()

			if p.Status != Working{
				log.Printf("ProxyName [%s] is not working, new user conn close", p.Name)
				c.Close()
				p.Unlock()
				return
			}
			p.UserConnList.PushBack(c)
			p.Unlock()

			p.CtlMsgChan <- 1
		}
	}()

	//配对
	go func(){
		for {
			cliConn := <-p.CliConnChan
			p.Lock()
			element := p.UserConnList.Front()

			var userConn *conn.Conn
			if element != nil{
				userConn = element.Value.(*conn.Conn)
				p.UserConnList.Remove(element)
			}else{
				cliConn.Close()
				continue
			}
			p.Unlock()

			log.Printf("Join two conns, (l[%s] r[%s]) (l[%s] r[%s])", cliConn.GetLocalAddr(), cliConn.GetRemoteAddr(),
				userConn.GetLocalAddr(), userConn.GetRemoteAddr())

			go conn.Join(cliConn, userConn)
		}
	}()

	return nil
}

func (p *ProxyServer) Close(){
	p.Lock()
	p.Status = Idle
	p.CtlMsgChan = make(chan int64)
	p.CliConnChan = make(chan *conn.Conn)
	p.UserConnList = list.New()
	p.Unlock()
}

func(p *ProxyServer)WaitUserConn(res int64){
	res = <-p.CtlMsgChan
	return
}
