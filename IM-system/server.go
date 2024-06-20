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
	fmt.Println("连接建立成功,别丢份儿啊！",
		"田文镜，我草泥马")
	//buffer := make([]byte, 1024)
	//用户上线，将用户加入到onlineMap中
	user := Newuser(conn, this)

	user.Online()

	Islive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			msg := string(buf[:n-1])

			user.DoMessage(msg)
			Islive <- true
		}
	}()
	for {
		select {
		case <-Islive:
			// 阻塞
		case <-time.After(time.Second * 200):
			// 阻塞
			user.SendMsg("你被踢了")
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()
	i := 0
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		go this.process(conn)
		i++
		fmt.Printf("第%v次连接\n", i)
	}
}
