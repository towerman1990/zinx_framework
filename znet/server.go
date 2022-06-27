package znet

import (
	"fmt"
	"net"
	"zinx_framework/conf"
	"zinx_framework/ziface"
)

type Server struct {
	Name string

	IPVersion string

	IP string

	Port int

	MsgHandler ziface.IMsgHandle

	ConnManager ziface.IConnManager

	OnConnStart func(conn ziface.IConnection)

	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d, is starting\n",
		conf.Config.Name, conf.Config.Host, conf.Config.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		conf.Config.Version, conf.Config.MaxConn, conf.Config.MaxPackageSize)
	fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		s.MsgHandler.StartWorkPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resovle tcp addr error: ", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("resovle tcp addr error: ", err)
			return
		}

		fmt.Println("start zinx server success,", s.Name, "listenning...")
		var cid uint32 = 0
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			if s.ConnManager.Len() > conf.Config.MaxConn {
				fmt.Println("Too Many Connections, MaxConn =", conf.Config.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Zinx server name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router success")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFnc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFnc
}

func (s *Server) SetOnConnStop(hookFnc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFnc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        conf.Config.Name,
		IPVersion:   "tcp4",
		IP:          conf.Config.Host,
		Port:        conf.Config.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}
