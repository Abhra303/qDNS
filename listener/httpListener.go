package listener

import (
	"fmt"
	"net"
)

var DefaultPort int = 53

func PortListener(port int) {
	if port < 0 || port > 65353 {
		fmt.Println("the given port is invalid. using the default port 53...")
		port = DefaultPort
	}

	udpAddr := net.UDPAddr{
		Port: port,
	}
	udpConn, err := net.ListenUDP("udp", &udpAddr)

	if err != nil {
		panic(err)
	}

	for {
		inputBytes := make([]byte, 512)

		length, clientAddr, err := udpConn.ReadFromUDP(inputBytes)
		if err != nil {
			fmt.Println("error reading UDP packet")
		}

		fmt.Printf("clientAddr.Zone: %s\n", clientAddr.Zone)
		fmt.Printf("clientAddr.Network: %s\n", clientAddr.Network())
		fmt.Printf("clientAddr.ToString: %s\n", clientAddr.String())
		fmt.Printf("clientAddr.IP: %s\n", clientAddr.IP.String())
		fmt.Println("data: ", inputBytes[:length])
		fmt.Println()

		// go resolveDNSRequest(clientAddr)

	}

}
