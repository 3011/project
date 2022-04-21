package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	TimeOut    = 3 * time.Second
	RetryCount = 3
	wg         sync.WaitGroup
)

func main() {
	flagUDP := flag.Bool("udp", false, "Scan UDP. (default false)")
	// flagTCP := flag.Bool("tcp", false, "Just scan TCP. (default false)")
	// flagAll := flag.Bool("all", true, "Scan TCP and TCP. (default true)")
	flagPort := flag.String("port", "0,65535", "Port range.")
	flag.Parse()

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

	if *flagUDP {
		for i := portStart; i <= portStop; i++ {
			wg.Add(1)
			go scanAll(host, strconv.Itoa(i))
		}
	} else {
		for i := portStart; i <= portStop; i++ {
			wg.Add(1)
			go scanTCP(host, strconv.Itoa(i))
		}
	}

	wg.Wait()
}

func scanAll(host string, port string) {
	defer wg.Done()
	status := ""
	for i := 0; i < RetryCount; i++ {
		_, err := net.DialTimeout("tcp", host+":"+port, TimeOut)
		if err == nil {
			status += "TCP "
			break
		}
	}

	for i := 0; i < RetryCount; i++ {
		_, err := net.DialTimeout("udp", host+":"+port, TimeOut)
		if err == nil {
			status += "UDP"
			break
		}
	}
	if status != "" {
		fmt.Printf("%s %5s %s\n", host, port, status)
	}
}

func scanTCP(host string, port string) {
	defer wg.Done()
	for i := 0; i < RetryCount; i++ {
		_, err := net.DialTimeout("tcp", host+":"+port, TimeOut)
		if err == nil {
			fmt.Printf("%s %5s %s\n", host, port, "TCP")
			break
		}
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
