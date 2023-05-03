package resolver

import (
	"fmt"
	"net"

	"github.com/abhra303/qDNS/dnsparser"
	"github.com/abhra303/qDNS/zonefiles"
)

func ResolveDNSRequest(inputBytes []byte, length int, conn *net.UDPConn, clientAddr *net.UDPAddr) {
	query, err := dnsparser.ParseDnsQuery(inputBytes, length)
	if err != nil {
		fmt.Print(err)
		return
	}

	if !query.Header.QR {
		err = fmt.Errorf("corrupted header: message received as a response")
		fmt.Print(err)
		return
	}
	if query.Header.Opcode == 0 && query.Question == nil {
		err = fmt.Errorf("corrupted message: standard query don't have a question section")
		fmt.Print(err)
		return
	} else if query.Header.Opcode == 1 {
		fmt.Print("Inverse Query not supported for now\n")
		return
	} else if query.Header.Opcode == 2 {
		fmt.Print("Not supported yet")
		return
	}
	rrQuery := zonefiles.QueryDomain{QdCount: query.Header.Qdcount, Questions: query.Question}

	rrResults, err := zonefiles.SearchResourceRecords(&rrQuery)
	if err != nil {
		fmt.Print(err)
		return
	}

	response := dnsparser.DnsMessage{}
	response.Header = query.Header
	response.Header.Ancount = rrResults.Ancount
	response.Header.Arcount = rrResults.Arcount
	response.Header.Nscount = rrResults.Nscount
	response.Header.Rcode = rrResults.RCode
	response.Header.Z = 0
	response.Header.RA = true
	response.Header.QR = true
	response.Answer = rrResults.Answers
	response.Authority = rrResults.Authority
	response.Additional = rrResults.Additional

	rawMessage, err := dnsparser.SerializeMessage(&response)
	if err != nil {
		fmt.Print(err)
		return
	}
	_, err = conn.WriteToUDP(rawMessage, clientAddr)
	if err != nil {
		fmt.Printf("can't send message to client\n")
	}
	/* send the response back to the client/resolver */
	fmt.Println(query)
	fmt.Println(rrResults)
	fmt.Println(rawMessage)
}
