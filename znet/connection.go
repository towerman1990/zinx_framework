package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx_framework/ziface"
)

type Connection struct {
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	ExitChan chan bool

	msgChan chan []byte

	MsgHandler ziface.IMsgHandle
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID =", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		go c.MsgHandler.DoMsgHandler(&req)
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID =", c.ConnID)

	go c.StartReader()

	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn.Stop... ConnID =", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgID)
		return errors.New("pack error msg")
	}

	c.msgChan <- binaryMsg

	return nil
}
