package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func Newuser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
func (this *User) Online() {
	this.server.maplock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.maplock.Unlock()
	this.server.BroadCast(this, "已上线")
}
func (this *User) Offline() {
	this.server.maplock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.maplock.Unlock()
	this.server.BroadCast(this, "已下线")
}
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.maplock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.maplock.Unlock()
		return
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		this.server.maplock.Lock()
		newname := msg[7:]
		delete(this.server.OnlineMap, this.Name)
		this.server.OnlineMap[newname] = this
		this.Name = newname
		this.server.maplock.Unlock()
		this.SendMsg("已修改用户名\n")
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remotename := strings.Split(msg, "|")[1]

		if remotename == "" {
			this.SendMsg("消息格式不正确，请使用\"to|张三|你好\"。\n")
			return
		}
		if remotename == this.Name {
			this.SendMsg("不能给自己发送消息。\n")
			return
		}

		remoteuser, ok := this.server.OnlineMap[remotename]
		if !ok {
			this.SendMsg("该用户不存在。\n")
			return
		}

		content := strings.Split(msg, "|")[2]
		remoteuser.SendMsg(this.Name + "对你说:" + content)
	} else if msg == "me" {
		this.SendMsg("我的用户名是:\n" + this.Name + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}

}
