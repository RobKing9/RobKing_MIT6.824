package worker

import (
	"6.824/src/mr"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

//做 Reduce 任务
// 首先需要获得这些 中间文件的文件名
func (w *worker) doReduceTask(t mr.Task) {
	log.Println("准备开始做Reduce任务！")
	//存放所有文件的 内容
	hashTable := make(map[string][]string)
	//获取中间文件 文件名
	for i := 0; i < t.NMap; i++ {
		reduceFileName := reduceFileName(i, t.TaskId)
		//打开这些文件
		file, err := os.Open(reduceFileName)
		if err != nil {
			w.ReportTask(t, false, err)
			return
		}
		//解码文件内容
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			//文件解码结束
			if err := dec.Decode(&kv); err != nil {
				break
			}
			//没有键值
			if _, ok := hashTable[kv.Key]; !ok {
				hashTable[kv.Key] = make([]string, 0)
			}
			hashTable[kv.Key] = append(hashTable[kv.Key], kv.Value)
		}
		os.Remove(reduceFileName)
	}
	//保存 reduce 任务的结果
	res := make([]string, 0)
	for k, v := range hashTable {
		res = append(res, fmt.Sprintf("%s %s\n", k, w.reducef(k, v)))
	}
	//将结果写入到 文件中
	file, err := os.Create(fmt.Sprintf("mr-out-%d", t.TaskId))
	defer file.Close()
	if err != nil {
		w.ReportTask(t, false, err)
		return
	}
	_, err = file.Write([]byte(strings.Join(res, "")))
	if err != nil {
		w.ReportTask(t, false, err)
		return
	}
	log.Println("你成功完成了Reduce任务，太棒了！！！")
	w.ReportTask(t, true, nil)
}
