package rpcReq

// RPC 参数和返回

import (
	"6.824/src/mr"
)

// ReqWorkerIdArgs 申请工号的参数
type ReqWorkerIdArgs struct {
}

// ReqWorkerIdReply 返回的 工号
type ReqWorkerIdReply struct {
	WorkerId int
}

// TaskArgs 申请任务的参数
type TaskArgs struct {
	WorkerId int
}

// TaskReply 返回的任务
type TaskReply struct {
	Task *mr.Task
}

// ReportTaskArgs 汇报任务
type ReportTaskArgs struct {
	Done      bool         //是否完成
	WorkId    int          //工人号
	TaskId    int          //任务号
	TaskPhase mr.TaskPhase //任务类型
}

type ReportTaskReply struct {
}
