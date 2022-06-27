package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

const protocol = "tcp"
const addr = "127.0.0.1:8088"

func TestDataPack(t *testing.T) {
	listenner, err := net.Listen(protocol, addr)
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept err:")
			}

			go func(conn net.Conn) {

				dp := NewDataPack()

				for {
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head error:", err)
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("server unpack data error:", err)
							return
						}

						fmt.Println("-> Recv MsgID:", msg.ID, ", DataLen: ", msg.DataLen, ", Data: ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Println("client dial error:", err)
		return
	}

	dp := NewDataPack()

	msg1 := &Message{
		ID:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &Message{
		ID:      2,
		DataLen: 6,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	select {}
}
