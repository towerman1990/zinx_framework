package znet

import (
	"fmt"
	"strconv"
	"zinx_framework/conf"
	"zinx_framework/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter

	TaskQueue []chan ziface.IRequest

	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: conf.Config.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, conf.Config.MaxWorkerTaskLen),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), "is NOT FOUND! Need Register!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID =" + strconv.Itoa(int(msgID)))
	}

	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID, " success")
}

func (mh *MsgHandle) StartWorkPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, conf.Config.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID =", workerID, " is started...")
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() & mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		", request MsgID = ", request.GetMsgID(),
		", workerID = ", workerID)

	mh.TaskQueue[workerID] <- request
}
