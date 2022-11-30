package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abhra303/qDNS/listener"
)

func main() {
	var port int
	var err error
	arguments := os.Args

	if len(arguments) == 2 {
		port, err = strconv.Atoi(arguments[1])

		if err != nil {
			fmt.Println("the given port number is not an integer:", arguments[1])
			return
		}
	} else {
		port = listener.DefaultPort
	}

	fmt.Printf("starting server at port %v ...\n", port)

	udpConn := listener.PortListener(port)

	// the infinite loop which looks for udp packets
	for {
		inputBytes := make([]byte, 512)

		length, clientAddr, err := udpConn.ReadFromUDP(inputBytes)
		if err != nil {
			fmt.Println("error reading UDP packet")
			continue
		}

		fmt.Printf("clientAddr.Zone: %s\n", clientAddr.Zone)
		fmt.Printf("clientAddr.Network: %s\n", clientAddr.Network())
		fmt.Printf("clientAddr.ToString: %s\n", clientAddr.String())
		fmt.Printf("clientAddr.IP: %s\n", clientAddr.IP.String())
		fmt.Println("data: ", inputBytes[:length])

		// go resolveDNSRequest(clientAddr)

	}
}
