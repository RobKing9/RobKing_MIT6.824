package worker

import (
	"6.824/src/mr"
	"6.824/src/mr/rpcReq"
	"log"
	"net/rpc"
	"os"
)

//注册工号
//申请任务
//汇报任务

//给老板发送RPC请求
func call(rpcName string, args interface{}, reply interface{}) bool {
	//通过套接字 访问RPC服务
	c, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("连接 master 失败：", err.Error())
	}
	defer c.Close()

	err = c.Call(rpcName, args, reply)
	if err != nil {
		//reading body gob: attempt to decode into a non-pointer
		//reading body gob: local interface type *interface {} can only be decoded from remote interface type; received concrete type ReqWorkerIdArgs = struct { }
		log.Println("调用失败！", err.Error())
		return false
	}
	return true
}

// ReqWorkerId 每个 worker 启动的时候都要向 master 注册一个 工号
func (w *worker) ReqWorkerId() {
	//参数 都是指针类型的
	//否则会报错 reading body gob: attempt to decode into a non-pointer
	args := &rpcReq.ReqWorkerIdArgs{}
	reply := &rpcReq.ReqWorkerIdReply{}

	if ok := call("Master.RegWorker", args, reply); !ok {
		log.Fatal("申请工号 失败！")
	}
	log.Printf("申请工号成功！你的工号是：%d, 接下来去申请任务把！！！", reply.WorkerId)
	w.workerId = reply.WorkerId
}

// ReqTask 根据自己的工号 去申请任务
func (w *worker) ReqTask() mr.Task {
	//请求参数 是自己的工号
	args := &rpcReq.TaskArgs{}
	args.WorkerId = w.workerId
	//结果
	reply := &rpcReq.TaskReply{}
	//通过RPC请求任务
	if ok := call("Master.GetOneTask", args, reply); !ok {
		//请求任务失败
		log.Println("申请任务失败！")
		os.Exit(1)
	}
	log.Printf("工号为%d的工人申请任务成功了！！！接下来准备干活吧！！！", w.workerId)
	return *reply.Task
}
