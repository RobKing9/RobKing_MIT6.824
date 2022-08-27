package worker

//做 Map 任务
import (
	"6.824/src/mr"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
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
	//编号是从0开始的
	log.Printf("你需要处理的文件是：%s, 你处理的Map任务编号是：%d", t.FileName, t.TaskId)
	content, err := ioutil.ReadFile(t.FileName)
	if err != nil {
		log.Println("ioutil.ReadFile failed:", err.Error())
		//汇报失败
		w.ReportTask(t, false, err)
	}
	//mapf
	kvs := w.mapf(t.FileName, string(content))
	//存放map任务的结果
	reduces := make([][]KeyValue, t.NReduce)
	for _, kv := range kvs {
		//哈希值，相同的哈希值肯定是一样的，对 NReduce 取余，得到这么多Reduce任务
		idx := ihash(kv.Key) % t.NReduce
		reduces[idx] = append(reduces[idx], kv)
	}
	//生成 很多 mr-mapIndex-reduceIndex 文件
	for i, reduce := range reduces {
		reduceFileName := reduceFileName(t.TaskId, i)
		//创建中间文件
		file, err := os.Create(reduceFileName)
		if err != nil {
			log.Println("os.create failed: ", err.Error())
			//汇报失败
			w.ReportTask(t, false, err)
			return
		}
		enc := json.NewEncoder(file)
		for _, kv := range reduce {
			//将内容编码 保存到文件中
			if err := enc.Encode(&kv); err != nil {
				log.Println("enc.Encode failed:", err.Error())
				//汇报失败
				w.ReportTask(t, false, err)
				return
			}
		}
		if err := file.Close(); err != nil {
			log.Println("file.Close failed:", err.Error())
			//汇报失败
			w.ReportTask(t, false, err)
			return
		}
	}
	log.Println("你成功完成年Map任务，太棒了！！！")
	//汇报完成
	w.ReportTask(t, true, nil)
}

//通过map任务的序号和reduce任务的序号，形成Map任务结果的中间文件的格式
func reduceFileName(mapId, reduceId int) string {
	return fmt.Sprintf("mr-%d-%d", mapId, reduceId)
}
