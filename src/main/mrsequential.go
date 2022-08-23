package main

//这是 非分布式的 实现

import (
	"6.824/src/mr/worker"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"sort"
)

// ByKey for sorting by key.
//结构体 数组
//主要是 为了 结构体排序
type ByKey []worker.KeyValue

func (a ByKey) Len() int      { return len(a) }
func (a ByKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less 按照 Key 进行升序排列
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main() {
	//两个参数 第一个是插件.so文件 第二个是电子书  len(os.Args)=3
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: mrsequential xxx.so inputfiles...\n")
		os.Exit(1)
	}
	//os.Args[1]是插件 .so文件
	mapf, reducef := loadPlugin(os.Args[1])

	//读取输入的文件 通过map任务化为键值对  累积map的输出
	intermediate := []worker.KeyValue{}
	//对于输入文件的处理
	//Args[2:] 代表的是输入的所有的文件
	for _, filename := range os.Args[2:] {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		//读取每一个文件的内容
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		//每一个文件都需要执行map任务 返回很多键值对 []mr.KeyValue
		kva := mapf(filename, string(content))
		intermediate = append(intermediate, kva...)
	}

	// a big difference from real MapReduce is that all the
	// intermediate data is in one place, intermediate[],
	// rather than being partitioned into NxM buckets.
	//按照 Key 升序排列
	sort.Sort(ByKey(intermediate))

	oname := "mr-out-0"
	//输出文件
	ofile, _ := os.Create(oname)

	//执行reduce任务 并将结果输出到 mr-out-* 文件
	i := 0
	for i < len(intermediate) {
		//j代表每个单词的个数
		j := i + 1
		//跳过相同的单词 使下一个 i  指向的是不同的
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		//相同的单词 出现的次数
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		//string 为每个单词出现的次数
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()
}

//参数是 .so文件
//返回的是两个函数 一个函数是mapf 执行map任务 另一个是reducef 执行reduce任务
func loadPlugin(filename string) (func(string, string) []worker.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	//调用Map方法
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	//mapf 就是 Map方法
	mapf := xmapf.(func(string, string) []worker.KeyValue)
	//调用Reduce方法
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
