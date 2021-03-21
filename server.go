package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//online user list
	OnlineMap map[string]*User //???
	mapLock   sync.RWMutex     //writer lock

	//message broadcast channel
	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User), //??
		Message:   make(chan string),
	}

	return server
}

//
func (this *Server) BroadCast(user *User, msg string) {
	fmt.Println("BroadCast")
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg //send sendMsg to chan Message

}

//对应一个客户端
func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	fmt.Println("链接建立成功")

	//...当前链接的业务
	//fmt.Println("链接建立成功")
	fmt.Println("Handler Handler")
	user := NewUser(conn, this)
	fmt.Printf("%+v\n", user)
	user.Online() //需执行msg := <-this.Message来解除阻塞
	fmt.Println("Handler islive1")
	//监听用户是否活跃的channel
	islive := make(chan bool)
	//接受客户端发送的消息
	// fmt.Println("Handler islive2")
	go func() {
		fmt.Println("Handler func")
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])
			fmt.Println("Handler sendMsg")
			//用户针对msg进行消息处理
			user.DoMessage(msg)
			islive <- true
		}
	}()

	//当前handler阻塞
	// select {} //select
	fmt.Println("当前handler阻塞")
	for {
		select {
		case <-islive: //do nothing,will renew status
		case <-time.After(time.Second * 100): //seconds
			//已经超时
			//将当前的User强制的关闭
			user.sendMsg("User强制的关闭")

			//销毁用的资源
			close(user.C) //??

			//关闭连接
			conn.Close() //?

			//退出当前Handler
			return //runtime.Goexit()
		}
	}

}

//启动服务器的接口
func (this *Server) Start() {
	//socket listen
	fmt.Println("链接")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept() //?? 当有客户端连接时执行下面
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}

//监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		// fmt.Println("listener msg msg:", msg)
		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg //SEND *user C
		}
		this.mapLock.Unlock()
	}
}
