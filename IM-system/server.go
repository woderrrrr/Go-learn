package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	maplock   sync.RWMutex
	//消息广播
	Msg chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Msg:       make(chan string),
	}
	return server
}

// 监听广播消息的读协程
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Msg
		this.maplock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.maplock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Msg <- sendMsg
}
func (this *Server) process(conn net.Conn) {
	defer conn.Close()
	fmt.Println("链接建立成功,田文镜，我草泥马")
	//buffer := make([]byte, 1024)
	//用户上线，将用户加入到onlineMap中
	user := Newuser(conn)

	this.maplock.Lock()
	this.OnlineMap[user.Name] = user
	this.maplock.Unlock()
	this.BroadCast(user, "已上线")

	select {}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		go this.process(conn)
	}
}
