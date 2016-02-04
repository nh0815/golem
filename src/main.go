package main

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	go poll()
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.Handle("/ws/", sockjs.NewHandler("/ws", sockjs.DefaultOptions, wsHandler))
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func wsHandler(session sockjs.Session) {
	log.Println("new sockjs session established")
}

func poll() {
	for {
		//log.Println("reading status")
		read_status()
		//log.Println("sleeping for 1s")
		time.Sleep(1000 * time.Millisecond)
	}
}

func read_status() [4444]string {
	var result [4444]string
	result[0] = read_cpu_info()
	result[1] = read_mem_info()
	result[2] = read_disk_info()
	result[3] = read_net_info()
	return result
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
