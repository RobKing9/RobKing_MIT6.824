package master

import (
	"6.824/src/mr"
	"sync"
	"time"
)

//分配的任务的状态
const (
	TaskReady = 0 // TaskReady 还没有分配的任务

	TaskChan = 1 //在管道中的任务

	TaskRunning = 2 // TaskRunning 分配了 正在运行的任务

	TaskFinish = 3 // TaskFinish 分配的已经完成了的

	TaskErr = 4 // TaskErr 执行过程出现错误
)

const (
	MaxTaskRunTime   = time.Second * 10       //任务运行的最长时间
	GenerateInterval = time.Millisecond * 500 //分配任务的间隔时间
)

// TaskStat 任务的动态信息
type TaskStat struct {
	status    int       //任务的状态
	workerId  int       //由 指定  worker
	startTime time.Time //开始做的时间
}

// Master 分配任务者
type Master struct {
	// Your definitions here.
	workerId  int          //发放的工号
	taskStats []TaskStat   //所有任务的动态信息
	taskPhase mr.TaskPhase //任务类型
	files     []string     //所有的文件
	nReduce   int          //reduce任务数量
	taskCh    chan mr.Task //发送任务的管道
	done      bool         //退出进程的标志 true意味着所有的任务都完成了
	mu        sync.Mutex
}

func MakeMaster(files []string, nReduce int) *Master {
	//初始化 master
	m := Master{}
	m.files = files
	m.nReduce = nReduce
	m.mu = sync.Mutex{}
	//管道的大小取决于 文件数量也就是Map任务数和Reduce任务数
	m.taskCh = make(chan mr.Task, max(len(m.files), m.nReduce))

	//一开始是做 Map任务 所以 master启动的时候就需要对Map任务进行初始化
	m.initMapTask()
	//管理任务，计划执行
	//开辟一个 携程 持续分配任务
	go func() {
		for !m.done {
			go m.GenerateTask()
			time.Sleep(GenerateInterval)
		}
	}()

	//与worker建立rpc通信
	//发放工号，分配任务
	m.server()
	return &m
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Done
//返回true证明所有的工作都完成了
func (m *Master) Done() bool {
	//ret := false
	return m.done
}
