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
	cpu := read_cpu_info()
	//println(strings.Split(cpu, "\n"))
	for _, element := range strings.Split(cpu, "\n") {
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
