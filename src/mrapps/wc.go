package main

import (
	"6.824/src/mr/worker"
	"strconv"
	"strings"
	"unicode"
)

// Map 第一个参数是电子书 第二个参数是取出来的内容 返回键值对
func Map(filename string, contents string) []worker.KeyValue {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	//将内容变为 单词数组 返回 []string
	words := strings.FieldsFunc(contents, ff)

	kva := []worker.KeyValue{}
	//每个单词的 value=1
	for _, w := range words {
		kv := worker.KeyValue{w, "1"}
		kva = append(kva, kv)
	}
	return kva
}

func Reduce(key string, values []string) string {
	// return the number of occurrences of this word.
	return strconv.Itoa(len(values))
}
