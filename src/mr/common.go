package mr

type TaskPhase int

const (
	MapPhase TaskPhase = 0 // MapPhase Map任务

	ReducePhase TaskPhase = 1 // ReducePhase Reduce任务
)

// Task 任务
type Task struct {
	TaskId int //任务序列号
	//任务类型
	TaskPhase TaskPhase
	FileName  string //map 任务 一个map任务 一个文件
	NReduce   int    //reduce任务数量
	NMap      int    //Map任务的数量
}
