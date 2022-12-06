package dnsparser

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

	Z       bool // for future use
	Rcode   int  // 4 bit response code
	Qdcount uint // no of entries in question section
	Ancount uint // no of RR in answer section
	Nscount uint // no of name server RR in authority records section
	Arcount uint // no of RR in the additional records section
}

type MessageQuestion struct {

	/*
	 A domain name represented as a sequence of labels, where
	 each label consists of a length octet followed by that
	 number of octets.  The domain name terminates with the
	 zero length octet for the null label of the root.  Note
	 that this field may be an odd number of octets; no
	 padding is used.
	*/
	QName string

	/*
	 A two octet code which specifies the type of the query.
	 The values for this field include all codes valid for a
	 TYPE field, together with some more general codes which
	 can match more than one type of RR.
	*/
	Qtype int

	/*
	 A two octet code that specifies the class of the query.
	 For example, the QCLASS field is IN for the Internet.
	*/
	Qclass int
}

type ResourceRecord struct {

	/*
	 A domain name to which this resource record pertains
	*/
	Name string

	/*
	 Two octets containing one of the RR type codes.  This
	 field specifies the meaning of the data in the RDATA
	 field.
	*/
	Type int

	/*
	 two octets which specify the class of the data in the
	 RDATA field.
	*/
	Class int

	/*
	 A 32 bit unsigned integer that specifies the time
	 interval (in seconds) that the resource record may be
	 cached before it should be discarded.  Zero values are
	 interpreted to mean that the RR can only be used for the
	 transaction in progress, and should not be cached.
	*/
	TTL uint

	/*
	 An unsigned 16 bit integer that specifies the length in
	 octets of the RDATA field.
	*/
	Rdlength uint

	/*
	 A variable length string of octets that describes the
	 resource.  The format of this information varies
	 according to the TYPE and CLASS of the resource record.
	 For example, the if the TYPE is A and the CLASS is IN,
	 the RDATA field is a 4 octet ARPA Internet address.
	*/
	Rdata string
}

type DnsMessage struct {
	Header     *MessageHeader
	Question   *MessageQuestion
	Answer     []*ResourceRecord
	Authority  []*ResourceRecord
	Additional []*ResourceRecord
}

func ParseDnsQuery(inputBytes []byte, length int) (DnsMessage, error) {

	return DnsMessage{}, nil
}
