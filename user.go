package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn //对端客户端
	server *Server
}

//create user

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//user channel 's goroutine
	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	//用户上线,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this //???   key:this.Name?  value:this(User)
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {

	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name) //get it
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")

}

func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

//监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C //this chan recieved
		this.conn.Write([]byte(msg + "\n"))
	}

}
