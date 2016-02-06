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
	Timestamp time.Time `json:"timestamp"`
	Cpu       CpuInfo   `json:"cpu"`
	Mem       MemInfo   `json:"memory"`
	Net       NetInfo   `json:"network"`
	Disk      DiskInfo  `json:"disk"`
}

type CpuInfo struct {
	User      int64 `json:"user"`
	Nice      int64 `json:"nice"`
	System    int64 `json:"system"`
	Idle      int64 `json:"idle"`
	Iowait    int64 `json:"iowait"`
	Irq       int64 `json:"irq"`
	IrqSoft   int64 `json:"irqSoft"`
	Steal     int64 `json:"steal"`
	Guest     int64 `json:"guest"`
	GuestNice int64 `json:""guestNice`
}

type MemInfo struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
}

type NetInfo struct {
	Interfaces []NetInterface `json:"interfaces"`
}

type NetInterface struct {
	Name               string `json:"name"`
	RecvBytes          int64  `json:"receiveBytes"`
	RecvPackets        int64  `json:"receivePackets"`
	RecvErrs           int64  `json:"receiveErrors"`
	RecvDrops          int64  `json:"receiveDrops"`
	RecvFifo           int64  `json:"receiveFifo"`
	RecvFrame          int64  `json:"receiveFrame"`
	RecvCompressed     int64  `json:"receiveCompressed"`
	RecvMulticast      int64  `json:"receiveMulticast"`
	TransmitBytes      int64  `json:"transmitBytes"`
	TransmitPackets    int64  `json:"transmitPackets"`
	TransmitErrs       int64  `json:"transmitErrors"`
	TransmitDrops      int64  `json:"transmitDrops"`
	TransmitFifo       int64  `json:"transmitFifo"`
	TransmitCollisions int64  `json:"transmitCollisions"`
	TransmitCarrier    int64  `json:"transmitCarrier"`
	TransmitCompressed int64  `json:"transmitCompressed"`
}

type DiskInfo struct {
	Disks []Disk `json:"disks"`
}

type Disk struct {
	Name            string `json:"name"`
	ReadsCompleted  int64  `json:"readsCompleted"`
	ReadsMerged     int64  `json:"readsMerged"`
	SectorsRead     int64  `json:"sectorsRead"`
	TimeReading     int64  `json:"timeReading"`
	WritesCompleted int64  `json:"writesCompleted"`
	WritesMerged    int64  `json:"writesMerged"`
	SectorsWritten  int64  `json:"sectorsWritten"`
	TimeWriting     int64  `json:"timeWriting"`
	IopsInProgress  int64  `json:"iopsInProgress"`
	IOTime          int64  `json:"ioTime"`
	IOTimeWeighted  int64  `json:"ioTimeWeighted"`
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
	system := SysInfo{time.Now(), read_cpu_info(), read_mem_info(), read_net_info(), read_disk_info()}
	return system
}

func read_cpu_info() CpuInfo {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}
	cpu := split_on_newline(string(data))[0]
	fields := strings.Fields(cpu)
	user := string_to_int64(fields[1])
	nice := string_to_int64(fields[2])
	system := string_to_int64(fields[3])
	idle := string_to_int64(fields[4])
	iowait := string_to_int64(fields[5])
	irq := string_to_int64(fields[6])
	irq_soft := string_to_int64(fields[7])
	steal := string_to_int64(fields[8])
	guest := string_to_int64(fields[9])
	guest_nice := string_to_int64(fields[10])
	cpu_info := CpuInfo{user, nice, system, idle, iowait, irq, irq_soft, steal, guest, guest_nice}
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

func read_disk_info() DiskInfo {
	disks := []Disk{}
	data, err := ioutil.ReadFile("/proc/diskstats")
	if err != nil {
		panic(err)
	}
	data_strings := split_on_newline(strings.TrimSpace(string(data)))
	for i := range data_strings {
		disk_info := strings.Fields(data_strings[i])[2:]
		name := disk_info[0]
		reads_completed := string_to_int64(disk_info[1])
		reads_merged := string_to_int64(disk_info[2])
		sectors_read := string_to_int64(disk_info[3])
		time_reading := string_to_int64(disk_info[4])
		writes_completed := string_to_int64(disk_info[5])
		writes_merged := string_to_int64(disk_info[6])
		sectors_written := string_to_int64(disk_info[7])
		time_writing := string_to_int64(disk_info[8])
		io_progress := string_to_int64(disk_info[9])
		io_time := string_to_int64(disk_info[10])
		io_time_weighted := string_to_int64(disk_info[11])
		disk := Disk{name, reads_completed, reads_merged, sectors_read, time_reading, writes_completed, writes_merged, sectors_written, time_writing, io_progress, io_time, io_time_weighted}
		disks = append(disks, disk)
	}
	return DiskInfo{disks}
}

func read_net_info() NetInfo {
	data, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		panic(err)
	}
	data_string := split_on_newline(strings.TrimSpace(string(data)))[2:]
	interfaces := []NetInterface{}
	for i := range data_string {
		interface_string := strings.Fields(data_string[i])
		name := strings.Replace(interface_string[0], ":", "", -1)
		recv_bytes := string_to_int64(interface_string[1])
		recv_packets := string_to_int64(interface_string[2])
		recv_errors := string_to_int64(interface_string[3])
		recv_drop := string_to_int64(interface_string[4])
		recv_fifo := string_to_int64(interface_string[5])
		recv_frame := string_to_int64(interface_string[6])
		recv_compressed := string_to_int64(interface_string[7])
		recv_multicast := string_to_int64(interface_string[8])
		transmit_bytes := string_to_int64(interface_string[9])
		transmit_packets := string_to_int64(interface_string[10])
		transmit_errors := string_to_int64(interface_string[11])
		transmit_drops := string_to_int64(interface_string[12])
		transmit_fifo := string_to_int64(interface_string[13])
		transmit_collision := string_to_int64(interface_string[14])
		transmit_carrier := string_to_int64(interface_string[15])
		transmit_compressed := string_to_int64(interface_string[16])
		net_interface := NetInterface{name, recv_bytes, recv_packets, recv_errors, recv_drop, recv_fifo, recv_frame, recv_compressed, recv_multicast, transmit_bytes, transmit_packets, transmit_errors, transmit_drops, transmit_fifo, transmit_collision, transmit_carrier, transmit_compressed}
		interfaces = append(interfaces, net_interface)
	}
	return NetInfo{interfaces}
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
