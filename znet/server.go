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
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d, is starting",
		conf.Config.Name, conf.Config.Host, conf.Config.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d",
		conf.Config.Version, conf.Config.MaxConn, conf.Config.MaxPackageSize)
	fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
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

			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router success")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       conf.Config.Name,
		IPVersion:  "tcp4",
		IP:         conf.Config.Host,
		Port:       conf.Config.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
