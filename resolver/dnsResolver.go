package resolver

import (
	"fmt"
	"net"

	"github.com/abhra303/qDNS/dnsparser"
)

func ResolveDNSRequest(inputBytes []byte, length int, clientAddr *net.UDPAddr) {
	message, err := dnsparser.ParseDnsQuery(inputBytes, length)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(message)
}