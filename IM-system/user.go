package main

import "net"

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
		this.server.OnlineMap[msg[7:]] = this
		delete(this.server.OnlineMap, this.Name)
		this.Name = msg[7:]
		this.server.maplock.Unlock()
		this.SendMsg("已修改用户名")
	} else {
		this.server.BroadCast(this, msg)
	}

}
