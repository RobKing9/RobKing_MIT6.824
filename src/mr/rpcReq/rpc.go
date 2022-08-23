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
