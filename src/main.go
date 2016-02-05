package main

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type SysInfo struct {
	Cpu CpuInfo `json:"cpu"`
}

type CpuInfo struct {
	User    string `json:"user"`
	Nice    string `json:"nice"`
	System  string `json:"system"`
	Idle    string `json:"idle"`
	Iowait  string `json:"iowait"`
	Irq     string `json:"irq"`
	Softirq string `json:"softirq"`
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
	system := SysInfo{read_cpu_info()}
	return system
}

func read_cpu_info() CpuInfo {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}
	cpu := split_on_newline(string(data))[0]
	fields := strings.Split(cpu, " ")
	cpu_info := CpuInfo{fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6]}
	return cpu_info
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
