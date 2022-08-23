package master

import (
	"6.824/src/mr"
	"log"
)

//对任务进行初始化

//初始化Map任务
func (m *Master) initMapTask() {
	//任务类型
	m.taskPhase = mr.MapPhase
	//所有 Map 任务的状态
	m.taskStats = make([]TaskStat, len(m.files))
}

//通过任务id 初始化一个任务
func (m *Master) initTaskById(taskId int) mr.Task {
	task := mr.Task{
		TaskId:    taskId,
		FileName:  "", //只有map任务 才有文件名
		TaskPhase: m.taskPhase,
	}
	if m.taskPhase == mr.MapPhase {
		task.FileName = m.files[taskId]
	}
	log.Printf("通过id初始化任务成功！任务号：%d，任务类型：%v，分配的任务文件是：%v", task.TaskId, task.TaskPhase, task.FileName)
	return task
}
