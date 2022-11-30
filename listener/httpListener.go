package listener

import (
	"fmt"
	"net"
)

var DefaultPort int = 53

func PortListener(port int) *net.UDPConn {
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

	return udpConn
}
