package main

import (
	//"bufio"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//"os"
)

func main() {
	//println("hello world")
	//info := read_cpu_info()
	//info := read_mem_info()
	//info := read_disk_info()
	/*
		info := read_net_info()
		for _, element := range split_on_newline(info) {
			println(element)
			println()
		}
	*/
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
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

func read_disk_info() string {
	data, err := ioutil.ReadFile("/proc/diskstats")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func read_net_info() string {
	data, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func split_on_newline(str string) []string {
	return strings.Split(str, "\n")
}
