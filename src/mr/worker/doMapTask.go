package worker

//做 Map 任务
import (
	"6.824/src/mr"
	"hash/fnv"
	"log"
)

//生成 Key 的哈希值
//相同的 Key 哈希值相同
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func (w *worker) doMapTask(t mr.Task) {
	log.Println("准备开始做 Map 任务！")
	log.Println("你需要处理的文件是：", t.FileName)
}
