package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abhra303/qDNS/config"
	"github.com/abhra303/qDNS/listener"
	"github.com/abhra303/qDNS/resolver"
	"github.com/abhra303/qDNS/zonefiles"
)

func main() {
	var path string
	var port int
	var err error
	arguments := os.Args

	if len(arguments) >= 2 {
		path = arguments[1]
	}
	if len(arguments) == 3 {
		port, err = strconv.Atoi(arguments[2])

		if err != nil {
			fmt.Println("the given port number is not an integer:", arguments[2])
			return
		}
	} else {
		port = listener.DefaultPort
	}

	err = config.LoadConfigFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !zonefiles.LoadZones() {
		fmt.Println("unable to load zones...")
		return
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

		go resolver.ResolveDNSRequest(inputBytes, length, udpConn, clientAddr)

		fmt.Printf("clientAddr.Zone: %s\n", clientAddr.Zone)
		fmt.Printf("clientAddr.Network: %s\n", clientAddr.Network())
		fmt.Printf("clientAddr.ToString: %s\n", clientAddr.String())
		fmt.Printf("clientAddr.IP: %s\n", clientAddr.IP.String())
		fmt.Println("data: ", inputBytes[:length])
	}
}
