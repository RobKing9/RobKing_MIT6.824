package worker

import (
	"6.824/src/mr"
	"time"
)

type worker struct {
	workerId int                             //每个工人都有自己的工号，表示自己做的任务
	mapf     func(string, string) []KeyValue //做map任务的能力
	reducef  func(string, []string) string   //做reduce任务的能力
}

// KeyValue
//键值对 key为单词 value为出现的次数
type KeyValue struct {
	Key   string
	Value string
}

// Worker
//参数 是map任务 和 Reduce任务
//通过RPC调用 分配任务的函数
func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	//初始化 worker
	w := worker{}
	w.mapf = mapf
	w.reducef = reducef
	//申请工号
	w.ReqWorkerId()
	//不断的申请任务
	for {
		task := w.ReqTask()
		w.doTask(task)
		time.Sleep(time.Second * 10)
	}
}

//做任务
func (w *worker) doTask(t mr.Task) {
	if t.TaskPhase == mr.MapPhase {
		w.doMapTask(t)
	}
}
