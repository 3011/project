package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	mothed int
	host   string
	port   string
}

var (
	Timeout      = time.Second * 5
	RetryCount   = 2
	MaxGoroutine int
	GoroutineNum = 0
	RequestChan  = make(chan Request)
	WorkerDone   = make(chan bool)
)

func main() {
	flagUDP := flag.Bool("udp", false, "是否扫描UDP 默认：否")
	// flagTCP := flag.Bool("tcp", false, "Just scan TCP. (default false)")
	// flagAll := flag.Bool("all", true, "Scan TCP and TCP. (default true)")
	maxGoroutine := flag.Int("maxGoroutine", 30, "最大协程数量 默认：30")
	flagPort := flag.String("port", "0,65535", "端口范围 默认：0,65535")
	flag.Parse()

	if *maxGoroutine < 0 || *maxGoroutine > 8000 {
		fmt.Println("GoroutineNum Error.")
	}
	MaxGoroutine = *maxGoroutine
	host := flag.Arg(0)
	if host == "" {
		fmt.Println("Please input host.")
		return
	}

	splitFlagPort := strings.Split(*flagPort, ",")
	var portStart int
	var portStop int
	if len(splitFlagPort) == 2 {
		var err error
		portStart, err = strconv.Atoi(splitFlagPort[0])
		if err == nil {
			portStop, err = strconv.Atoi(splitFlagPort[1])
		}
		if err != nil || (portStart > portStop) {
			fmt.Println("Error: -port " + *flagPort)
			return
		}
	} else {
		fmt.Println("Error: -port " + *flagPort)
		return
	}

	// if *flagTCP && !(*flagUDP) { // scan tcp
	// 	for i := portStart; i <= portStop; i++ {
	// 		wg.Add(1)
	// 		go scanTCP(host, strconv.Itoa(i))
	// 	}
	// } else if *flagUDP && !(*flagTCP) { //scan udp
	// 	for i := portStart; i <= portStop; i++ {
	// 		wg.Add(1)
	// 		go scanUDP(host, strconv.Itoa(i))
	// 	}
	// } else {
	// 	for i := portStart; i <= portStop; i++ {
	// 		wg.Add(1)
	// 		go scanAll(host, strconv.Itoa(i))
	// 	}
	// }
	GoroutineNum++
	go func() {
		if *flagUDP {
			for i := portStart; i <= portStop; i++ {
				if GoroutineNum < MaxGoroutine {
					RequestChan <- Request{0, host, strconv.Itoa(i)}
				} else {
					scanAll(host, strconv.Itoa(i), false)
				}
			}
			WorkerDone <- true
		} else {
			for i := portStart; i <= portStop; i++ {
				if GoroutineNum < MaxGoroutine {
					RequestChan <- Request{1, host, strconv.Itoa(i)}
				} else {
					scanTCP(host, strconv.Itoa(i), false)
				}
			}
			WorkerDone <- true
		}
	}()
	waitForWorkers()
}

func waitForWorkers() {
	for {
		select {
		case request := <-RequestChan:
			GoroutineNum++
			if request.mothed == 0 {
				go scanAll(request.host, request.port, true)
			} else if request.mothed == 1 {
				go scanTCP(request.host, request.port, true)
			}

		case <-WorkerDone:
			GoroutineNum--
			if GoroutineNum == 0 {
				return
			}
		}
	}

}

func scanAll(host string, port string, master bool) {
	status := ""
	for i := 0; i < RetryCount; i++ {

		_, err := net.DialTimeout("tcp", host+":"+port, Timeout)
		if err == nil {
			status += "TCP "
			break
		}
	}

	for i := 0; i < RetryCount; i++ {
		_, err := net.DialTimeout("udp", host+":"+port, Timeout)
		if err == nil {
			status += "UDP"
			break
		}
	}
	if status != "" {
		fmt.Printf("%s %5s %s\n", host, port, status)
	}
	if master {
		WorkerDone <- true
	}
}

func scanTCP(host string, port string, master bool) {
	for i := 0; i < RetryCount; i++ {
		// fmt.Printf("port: %v\n", port)
		_, err := net.DialTimeout("tcp", host+":"+port, Timeout)
		if err == nil {
			fmt.Printf("%s %5s %s\n", host, port, "TCP")
			break
		}
	}
	if master {
		WorkerDone <- true
	}
}

// func scanUDP(host string, port string) {
// 	defer wg.Done()
// 	for i := 0; i < RetryCount; i++ {
// 		_, err := net.DialTimeout("udp", host+":"+port, TimeOut)
// 		if err == nil {
// 			fmt.Printf("%s %5s %s\n", host, port, "UDP")
// 			break
// 		}
// 	}
// }
