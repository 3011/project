package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var TimeOut = time.Second * 5

func main() {
	flagUDP := flag.Bool("udp", false, "Just scan UDP. (default false)")
	flagTCP := flag.Bool("tcp", false, "Just scan TCP. (default false)")
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

	if *flagTCP && !(*flagUDP) { // scan tcp
		for i := portStart; i <= portStop; i++ {
			go scanTCP(host, strconv.Itoa(i))
		}
	} else if *flagUDP && !(*flagTCP) { //scan udp
		for i := portStart; i <= portStop; i++ {
			go scanUDP(host, strconv.Itoa(i))
		}
	} else {
		for i := portStart; i <= portStop; i++ {
			go scanAll(host, strconv.Itoa(i))
		}
	}
}

func scanAll(host string, port string) {
	status := ""
	_, err := net.DialTimeout("tcp", host+":"+port, TimeOut)
	if err == nil {
		status += "TCP "
	}
	_, err = net.DialTimeout("udp", host+":"+port, TimeOut)
	if err == nil {
		status += "UDP"
	}
	if status != "" {
		fmt.Println(host, port, status)
	}
}

func scanTCP(host string, port string) {
	_, err := net.DialTimeout("tcp", host+":"+port, TimeOut)
	if err == nil {
		fmt.Println(host, port, "TCP")
	}
}

func scanUDP(host string, port string) {
	_, err := net.DialTimeout("udp", host+":"+port, TimeOut)
	if err == nil {
		fmt.Println(host, port, "UDP")
	}
}
