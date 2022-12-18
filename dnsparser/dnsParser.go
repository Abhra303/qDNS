package dnsparser

import (
	"encoding/binary"
	"fmt"

	"github.com/abhra303/qDNS/zonefiles"
)

var headerSize int = 96         // message header size
var MessageByteLimit uint = 512 // overall message size

type MessageHeader struct {

	/*
	   A 16 bit identifier assigned by the program that
	   generates any kind of query.  This identifier is copied
	   the corresponding reply and can be used by the requester
	   to match up replies to outstanding queries.
	*/
	ID int

	QR bool // query(0) or response(1)

	/*
	   A four bit field that specifies kind of query in this
	   message.  This value is set by the originator of a query
	   and copied into the response.  The values are:

	      0               a standard query (QUERY)

	      1               an inverse query (IQUERY)

	      2               a server status request (STATUS)

	      3-15            reserved for future use
	*/
	Opcode int

	/*
	   Authoritative Answer - this bit is valid in responses,
	   and specifies that the responding name server is an
	   authority for the domain name in question section.

	   Note that the contents of the answer section may have
	   multiple owner names because of aliases.  The AA bit
	   corresponds to the name which matches the query name, or
	   the first owner name in the answer section.
	*/
	AA bool

	/*
	   TrunCation - specifies that this message was truncated
	   due to length greater than that permitted on the
	   transmission channel.
	*/
	TC bool

	/*
	   Recursion Desired - this bit may be set in a query and
	   is copied into the response. If RD is set, it directs
	   the name server to pursue the query recursively.
	   Recursive query support is optional.
	*/
	RD bool

	/*
	   Recursion Available - this be is set or cleared in a
	   response, and denotes whether recursive query support is
	   available in the name server.
	*/
	RA bool

	Z       int  // for future use
	Rcode   int  // 4 bit response code
	Qdcount uint // no of entries in question section
	Ancount uint // no of RR in answer section
	Nscount uint // no of name server RR in authority records section
	Arcount uint // no of RR in the additional records section
}

type DnsQuery struct {
	Header *MessageHeader

	/*
		One of the below two fields must be nil. Normal
		queries contain a Question section, but Inverse
		Queries contain Answer section in dns query and
		expect dns message responses containing Question
		section
	*/
	Question *[]*zonefiles.QueryQuestion
	Answer   []*zonefiles.ResourceRecord
}

type DnsMessage struct {
	Header     *MessageHeader
	Question   *[]*zonefiles.QueryQuestion
	Answer     []*zonefiles.ResourceRecord
	Authority  []*zonefiles.ResourceRecord
	Additional []*zonefiles.ResourceRecord
}

func parseQueryHeader(inputBytes []byte, length int, bytesOffset *int) *MessageHeader {
	header := MessageHeader{}

	if length < headerSize {
		return nil
	}

	header.ID = (int(inputBytes[*bytesOffset]) << 8) ^ int(inputBytes[*bytesOffset+1])
	*bytesOffset += 2

	header.QR = (int(inputBytes[*bytesOffset]) >> 8) == 1
	header.Opcode = int((int8(inputBytes[*bytesOffset]) << 1) >> 4)
	header.AA = (inputBytes[*bytesOffset] & 0b00000100) == 1
	header.TC = (inputBytes[*bytesOffset] & 0b00000010) == 1
	header.RD = (inputBytes[*bytesOffset] & 0b00000001) == 1

	*bytesOffset++
	header.RA = (inputBytes[*bytesOffset] & 0b10000000) == 1
	header.Z = int(inputBytes[*bytesOffset] & 0b01110000)
	header.Rcode = int(inputBytes[*bytesOffset] & 0b00001111)

	*bytesOffset++
	header.Qdcount = (uint(inputBytes[*bytesOffset]) << 8) ^ uint(inputBytes[*bytesOffset+1])
	*bytesOffset++
	header.Ancount = (uint(inputBytes[*bytesOffset]) << 8) ^ uint(inputBytes[*bytesOffset+1])
	*bytesOffset++
	header.Nscount = (uint(inputBytes[*bytesOffset]) << 8) ^ uint(inputBytes[*bytesOffset+1])
	*bytesOffset++
	header.Arcount = (uint(inputBytes[*bytesOffset]) << 8) ^ uint(inputBytes[*bytesOffset+1])
	return &header
}

func parseQueryQuestion(inputBytes []byte, bytesOffset *int) *zonefiles.QueryQuestion {
	bufLen := len(inputBytes)
	question := zonefiles.QueryQuestion{}
	fixedQSize := 2 + 2 // each question size should atleast 4 bytes long (2 byte QType + 2 byte QClass)

	if bufLen-*bytesOffset <= fixedQSize {
		return nil
	}

	// though it seems O(n^2) but actually is O(n); n is the no. of bytes in Qname
	for i, length := 0, int(inputBytes[*bytesOffset]); length != 0; i++ {
		qName := ""
		if *bytesOffset+length+fixedQSize > bufLen {
			return nil
		}

		for ; length > 0; length-- {
			*bytesOffset++
			qName += string(inputBytes[*bytesOffset])
		}

		if qName == "" {
			return nil
		}
		question.QName[i] = qName
		*bytesOffset++
	}

	question.Qtype = (int(inputBytes[*bytesOffset]) << 8) ^ int(inputBytes[*bytesOffset+1])
	*bytesOffset += 2
	question.Qclass = (int(inputBytes[*bytesOffset]) << 8) ^ int(inputBytes[*bytesOffset+1])
	*bytesOffset += 2

	return &question
}

func parseQueryQuestions(inputBytes []byte, qdCount uint, bytesOffset *int) *[]*zonefiles.QueryQuestion {
	var messageQuestions []*zonefiles.QueryQuestion
	var i uint

	for i = 0; i < qdCount; i++ {
		messageQuestion := parseQueryQuestion(inputBytes, bytesOffset)

		if messageQuestion == nil {
			return nil
		}
		messageQuestions = append(messageQuestions, messageQuestion)
	}
	return &messageQuestions
}

func ParseDnsQuery(inputBytes []byte, length int) (*DnsQuery, error) {
	var err error
	bytesOffset := 0
	query := DnsQuery{}

	query.Header = parseQueryHeader(inputBytes, length, &bytesOffset)
	if query.Header == nil {
		err = fmt.Errorf("error parsing dns request header")
		return &query, err
	}

	query.Question = parseQueryQuestions(inputBytes, query.Header.Qdcount, &bytesOffset)
	if query.Question == nil {
		err = fmt.Errorf("error parsing dns request question section")
		return &query, err
	}

	return &query, nil
}

func serializeMessageHeader(header *MessageHeader, rawMessage []byte, offset *uint) (uint, error) {
	binary.BigEndian.PutUint16(rawMessage, uint16(header.ID))
	*offset += 2

	if header.QR {
		rawMessage[*offset] = 1 << 7
	}

	rawMessage[*offset] |= (byte(header.Opcode) << 3) & 0b01111000
	if header.AA {
		rawMessage[*offset] |= 0x4
	}

	tcPos := *offset
	if header.RD {
		rawMessage[*offset] |= 0x1
	}
	*offset++
	if header.RA {
		rawMessage[*offset] = 1 << 7
	}
	rawMessage[*offset] |= (byte(header.Z) << 4) & 0b01110000
	rawMessage[*offset] |= byte(header.Rcode) & 0b00001111
	*offset++

	binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(header.Qdcount))
	*offset += 2
	binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(header.Ancount))
	*offset += 2
	binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(header.Nscount))
	*offset += 2
	binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(header.Arcount))

	return tcPos, nil
}

func serializeMessageQuestion(questions *[]*zonefiles.QueryQuestion, rawMessage []byte, offset *uint) (bool, error) {
	for _, question := range *questions {
		labels := question.QName
		for _, label := range labels {
			len := len(label)

			rawMessage[*offset] = byte(len)
			*offset++
			if *offset+3 >= MessageByteLimit {
				return true, nil
			}

			for i := range label {
				rawMessage[*offset] = byte(label[i])
				*offset++
			}
		}
		// a null byte to denote the end of query domain
		*offset++

		binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(question.Qtype))
		*offset += 2

		binary.BigEndian.PutUint16(rawMessage[*offset:], uint16(question.Qclass))
		*offset += 2
	}
	return (*offset >= MessageByteLimit), nil
}

func serializeResourceRecords(RRs []*zonefiles.ResourceRecord, rawMessage []byte, offset *uint) (bool, error) {
	return false, nil
}

func SerializeMessage(message *DnsMessage) ([]byte, error) {
	var offset uint
	var tcPos uint
	var isTruncated bool
	rawMessage := make([]byte, 512)

	tcPos, err := serializeMessageHeader(message.Header, rawMessage, &offset)
	if err != nil {
		return nil, err
	}

	isTruncated, err = serializeMessageQuestion(message.Question, rawMessage, &offset)
	if err != nil {
		return nil, err
	} else if isTruncated {
		goto truncated
	}

	isTruncated, err = serializeResourceRecords(message.Answer, rawMessage, &offset)
	if err != nil {
		return nil, err
	} else if isTruncated {
		goto truncated
	}

	isTruncated, err = serializeResourceRecords(message.Answer, rawMessage, &offset)
	if err != nil {
		return nil, err
	} else if isTruncated {
		goto truncated
	}

	return rawMessage, err

truncated:
	rawMessage[tcPos] &= ^byte(0x2)
	return rawMessage, err
}
