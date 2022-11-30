package resolver

import (
	"fmt"
	"net"

	"github.com/abhra303/qDNS/dnsparser"
)

func ResolveDNSRequest(inputBytes []byte, length int, clientAddr *net.UDPAddr) {
	message, err := dnsparser.ParseDnsMessage(inputBytes, length)
	if err != nil {
		fmt.Print("the given udp packet do not contain a valid dns message\n")
		return
	}

	fmt.Println(message)
}
