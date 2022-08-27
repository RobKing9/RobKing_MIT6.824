package master

import (
	"6.824/src/mr/rpcReq"
	"log"
	"net/http"
	"net/rpc"
	"time"
)

//注册RPC服务
//worker 申请工号
//worker 申请任务
//worker 汇报任务

// Server 注册RPC服务
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	log.Println("服务器启动成功！！！正在监听9999端口......")
	go http.ListenAndServe("127.0.0.1:9999", nil)
}

// RegWorker 给工人 发放工号
func (m *Master) RegWorker(args *rpcReq.ReqWorkerIdArgs, reply *rpcReq.ReqWorkerIdReply) error {
	//防止Worker竞争id
	m.mu.Lock()
	defer m.mu.Unlock()
	//工号每次 +1
	m.workerId = m.workerId + 1
	reply.WorkerId = m.workerId
	log.Println("成功给工人发放工号！工号是：", m.workerId)
	return nil
}

// GetOneTask 给 worker 发送任务
func (m *Master) GetOneTask(args *rpcReq.TaskArgs, reply *rpcReq.TaskReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	//从任务管道中 取出 一个任务
	task := <-m.taskCh
	//返回给 worker
	reply.Task = &task
	//改变这个任务的动态信息
	//分配给的工人，此时分配的时间
	m.taskStats[task.TaskId].workerId = args.WorkerId
	m.taskStats[task.TaskId].startTime = time.Now()
	m.taskStats[task.TaskId].status = TaskRunning
	log.Printf("成功给工号为%d的工人分配任务！此时时间是：%v", args.WorkerId, time.Now())
	return nil
}

// AcceptReport 接收 Worker 汇报任务
func (m *Master) AcceptReport(args *rpcReq.ReportTaskArgs, reply *rpcReq.ReportTaskReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	//判断 任务类型以及工人id 和分配的时候是不是一致的
	if args.TaskPhase != m.taskPhase || args.WorkId != m.taskStats[args.TaskId].workerId {
		log.Println("分配不一致！")
		return nil
	}
	if args.Done {
		log.Printf("工号为%d的工人成功完成了任务类型为%d的任务号%d！", args.WorkId, args.TaskPhase, args.TaskId)
		//任务完成
		m.taskStats[args.TaskId].status = TaskFinish
	} else {
		log.Println("没有成功完成任务")
		//执行的过程出错了
		m.taskStats[args.TaskId].status = TaskErr
	}
	//继续调度
	go m.GenerateTask()
	return nil
}
