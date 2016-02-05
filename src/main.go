package main

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SysInfo struct {
	Cpu CpuInfo `json:"cpu"`
	Mem MemInfo `json:"memory"`
}

type CpuInfo struct {
	User   string `json:"user"`
	Nice   string `json:"nice"`
	System string `json:"system"`
	Idle   string `json:"idle"`
	Iowait string `json:"iowait"`
}

type MemInfo struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
}

var broadcaster pubsub.Publisher

func main() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.Handle("/ws/", sockjs.NewHandler("/ws", sockjs.DefaultOptions, wsHandler))
	log.Println("Listening...")
	go func() {
		poll()
	}()
	http.ListenAndServe(":3000", nil)
}

func wsHandler(session sockjs.Session) {
	log.Println("new sockjs session established")
	go func() {
		reader, _ := broadcaster.SubChannel(nil)
		for {
			status := <-reader
			if err := session.Send(status.(string)); err != nil {
				return
			}
		}
	}()
}

func poll() {
	for {
		status := read_status()
		system_status, err := json.Marshal(status)
		if err != nil {
			log.Println(err)
		}
		go func() {
			broadcaster.Publish(string(system_status))
		}()
		time.Sleep(time.Second)
	}
}

func read_status() SysInfo {
	system := SysInfo{read_cpu_info(), read_mem_info()}
	return system
}

func read_cpu_info() CpuInfo {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}
	cpu := split_on_newline(string(data))[0]
	fields := strings.Fields(cpu)
	cpu_info := CpuInfo{fields[1], fields[2], fields[3], fields[4], fields[5]}
	return cpu_info
}

func read_mem_info() MemInfo {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	mem_info := split_on_newline(string(data))
	total_string_array := strings.Fields(mem_info[0])
	total := byte_string_to_bits(total_string_array[1], total_string_array[2])
	free_string_array := strings.Fields(mem_info[1])
	free := byte_string_to_bits(free_string_array[1], free_string_array[2])
	return MemInfo{total, free}
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

func byte_string_to_bits(bytes string, suffix string) int64 {
	suffix = strings.ToUpper(suffix)
	switch {
	case suffix == "B":
		return string_to_int64(bytes)
	case suffix == "KB":
		return string_to_int64(bytes) * 1024
	case suffix == "MB":
		return string_to_int64(bytes) * 1024 * 1024
	case suffix == "GB":
		return string_to_int64(bytes) * 1024 * 1024 * 1024
	case suffix == "TB":
		return string_to_int64(bytes) * 1024 * 1024 * 1024 * 1024
	}
	return string_to_int64(bytes)
}

func string_to_int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return i
}

func split_on_newline(str string) []string {
	return strings.Split(str, "\n")
}
