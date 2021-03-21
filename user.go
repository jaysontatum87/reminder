package main

import (
	"fmt"
	"net"
	"strings"
)

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
	this.server.BroadCast(this, "ONLINE")
}

func (this *User) Offline() {

	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name) //get it
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "ONLINE")

}

//
func (this *User) sendMsg(msg string) {

	fmt.Println("sendMsg")
	this.conn.Write([]byte(msg))
}
func (this *User) DoMessage(msg string) {
	//Lock
	if msg == "test" {
		fmt.Println("test")
		this.server.mapLock.Lock()
		for k, v := range this.server.OnlineMap {
			onlineMsg := "key:" + k + "[" + v.Addr + "]" + v.Name + ":" + "ONLINE...\n"
			this.sendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newname := msg[7:] //???
		// //消息格式: rename|张三
		// newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newname]
		if ok {
			fmt.Println("this name (key) had already existed")
			this.sendMsg("this name (key) had already existed")
		} else {
			//delete and create
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, newname)
			this.server.OnlineMap[newname] = this
			this.Name = newname
			this.server.mapLock.Unlock()

			this.sendMsg("this new name:" + newname)

		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//1 获取对方的用户名
		name := strings.Split(msg, "|")[1]
		if name == "" {
			this.sendMsg("name null")
			return
		}

		//2 根据用户名 得到对方User对象
		value, ok := this.server.OnlineMap[name]
		//value :=this.server.OnlineMap[name] if value==""{}????
		if !ok {
			this.sendMsg("value null")
			return
		}

		//3 获取消息内容，通过对方的User对象将消息内容发送过去

		content := strings.Split(msg, "|")[2]
		//value = user
		value.sendMsg("from " + this.Name + " msg: " + content)
	} else {
		this.server.BroadCast(this, msg)
	}
}

//监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C //this chan recieved
		this.conn.Write([]byte(msg + "\n"))
	}

}
