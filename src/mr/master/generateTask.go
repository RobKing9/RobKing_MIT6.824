package master

//给 worker 分配任务

// GenerateTask 管理所有的任务
func (m *Master) GenerateTask() {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	//任务都完成了
	if m.done {
		return
	}
	//通过判断任务 处于的状态 做不一样的事情
	for taskid, task := range m.taskStats {
		//处于 准备中 我们就把这个任务放到管道中
		if task.status == TaskReady {
			//根据 任务id 初始化一个任务
			m.taskCh <- m.initTaskById(taskid)
			m.taskStats[taskid].status = TaskChan
		}
	}
}
