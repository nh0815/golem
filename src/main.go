package main

import (
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

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
		//status := read_status()
		status := "asdf"
		go func() {
			log.Println("sending status")
			broadcaster.Publish(status)
			log.Println("sent status")
		}()
		//channel <- read_status()[0]
		log.Println("sleeping for 1000ms")
		time.Sleep(time.Second)
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
