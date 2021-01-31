package conn

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Listener struct {
	Addr net.Addr
	Conns chan *Conn
}

//等待获取一个连接
func(l *Listener)GetConn()(conn *Conn){
	conn = <-l.Conns
	return
}

type Conn struct {
	TcpConn *net.TCPConn
	Reader *bufio.Reader
}



//连接服务器
func(c *Conn) ConnServer(host string, port int64)(err error){
	serverAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil{
		return err
	}
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil{
		return err
	}
	c.TcpConn = conn
	c.Reader = bufio.NewReader(c.TcpConn)
	return nil
}

func(c *Conn)GetRemoteAddr()(addr string){
	return c.TcpConn.RemoteAddr().String()
}

func(c *Conn)GetLocalAddr()(addr string){
	return c.TcpConn.LocalAddr().String()
}

func(c *Conn)ReadLine()(buf string, err error){
	buf, err = c.Reader.ReadString('\n')
	return
}

func(c *Conn)Write(content string)(err error){
	_, err = c.TcpConn.Write([]byte(content))
	return err
}

func(c *Conn)Close(){
	c.TcpConn.Close()
}

func Listen(bindAddr string, bindPort int64)(l *Listener, err error){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil{
		return nil, err
	}

	l = &Listener{
		Addr: listener.Addr(),
		Conns: make(chan *Conn),
	}

	go func(){
		for{
			conn, err := listener.AcceptTCP()
			if err != nil{
				log.Printf("Accept new tcp connection error, %v", err)
				continue
			}

			c := &Conn{
				TcpConn: conn,
			}
			c.Reader = bufio.NewReader(c.TcpConn)
			l.Conns <- c
		}
	}()
	return l, err
}

//阻塞到连接关闭
func Join(c1 *Conn, c2 *Conn){
	var wait sync.WaitGroup
	pipe := func(to *Conn, from *Conn){
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		_, err := io.Copy(to.TcpConn, from.TcpConn)
		if err != nil{
			log.Printf("join conns error %v\n", err)
		}
	}

	wait.Add(2)
	go pipe(c1, c2)
	go pipe(c1, c2)

	wait.Wait()
	return
}