package master

import (
	"6.824/src/mr"
	"log"
	"time"
)

//给 worker 分配任务

// GenerateTask 管理所有的任务
func (m *Master) GenerateTask() {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	//任务都完成了
	if m.done {
		return
	}
	//标志任务是否全部完成
	allFinish := true
	//通过判断任务 处于的状态 做不一样的事情
	for taskid, task := range m.taskStats {
		//处于 准备中 我们就把这个任务放到管道中
		if task.status == TaskReady {
			allFinish = false
			//根据 任务id 初始化一个任务
			m.taskCh <- m.initTaskById(taskid)
			m.taskStats[taskid].status = TaskChan
		}
		//在管道中
		if task.status == TaskChan {
			allFinish = false
			continue
		}
		//超时了 没有完成
		if task.status == TaskRunning {
			allFinish = false
			if time.Now().Sub(m.taskStats[taskid].startTime) > MaxTaskRunTime {
				log.Println("超时没有完成，正在进行重新分配")
				//重新进行分配
				m.taskCh <- m.initTaskById(taskid)
				m.taskStats[taskid].status = TaskChan
			}
		}
		if task.status == TaskFinish {
			continue
		}
		//出错了 重新返回管道中 继续被执行
		if task.status == TaskErr {
			allFinish = false
			log.Printf("任务号%d出错了,正在重新进行分配", taskid)
			m.taskCh <- m.initTaskById(taskid)
			m.taskStats[taskid].status = TaskChan
		}
	}
	//完成了
	if allFinish {
		if m.taskPhase == mr.MapPhase {
			//Map 任务全部完成了，开始做Reduce任务
			m.initReduceTask()
		} else {
			//所有任务都完成了
			m.done = true
		}
	}
}
