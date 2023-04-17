package listener

import (
	"net"
)

var DefaultPort int = 53

func PortListener(port int) *net.UDPConn {
	udpAddr := net.UDPAddr{
		Port: port,
	}
	udpConn, err := net.ListenUDP("udp", &udpAddr)

	if err != nil {
		panic(err)
	}

	return udpConn
}
