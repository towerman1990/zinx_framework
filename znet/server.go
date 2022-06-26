package znet

import (
	"fmt"
	"net"
	"zinx_framework/ziface"
)

type Server struct {
	Name string

	IPVersion string

	IP string

	Port int
}

func (s *Server) Start() {
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
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Recv buf err: ", err)
						continue
					}

					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write back buf err: ", err)
						continue
					}
				}
			}()
		}
	}()
}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.Start()
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp",
		IP:        "0.0.0.0",
		Port:      8088,
	}
	return s
}
