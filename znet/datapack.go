package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx_framework/conf"
	"zinx_framework/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if conf.Config.MaxPackageSize > 0 && msg.DataLen > conf.Config.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	return msg, nil
}
