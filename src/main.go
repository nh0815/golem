package main

import (
	//"bufio"
	//"fmt"
	"io/ioutil"
	"strings"
	//"os"
)

func main() {
	//println("hello world")
	//info := read_cpu_info()
	info := read_mem_info()
	for _, element := range split_on_newline(info) {
		println(element)
		println()
	}
}

func read_cpu_info() string {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func read_mem_info() string {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func split_on_newline(str string) []string {
	return strings.Split(str, "\n")
}
