package main

import (
	"fmt"
	"net"
	"os"
	"sort"

	"github.com/akamensky/argparse"
)

func worker(ip string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", ip, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	var portcap int

	parser := argparse.NewParser("globe", "port scanner")

	var ip *string = parser.String("i", "ip", &argparse.Options{Required: true, Help: "ip/domain to target"})
	var ports_to_scan *string = parser.String("p", "port", &argparse.Options{Required: false, Help: "ports to scan"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *ports_to_scan == "all" {
		portcap = 65535
	} else {
		portcap = 1024
	}

	fmt.Println("IP: " + *ip)

	ports := make(chan int, 100)
	results := make(chan int)

	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(*ip, ports, results)
	}

	go func() {
		for i := 1; i <= portcap; i++ {
			ports <- i
		}
	}()

	for i := 0; i < portcap; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	fmt.Printf(" PORT \t STATE \n======\t=======\n")
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf(" %d \t open \n", port)
	}
}
