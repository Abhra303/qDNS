package zonefiles

type QueryQuestion struct {

	/*
	   A domain name represented as a sequence of labels, where
	   each label consists of a length octet followed by that
	   number of octets.  The domain name terminates with the
	   zero length octet for the null label of the root.  Note
	   that this field may be an odd number of octets; no
	   padding is used.
	*/
	QName []string

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

/*
QueryDomain contains the required informations to perform
domain query.
*/
type QueryDomain struct {
	QdCount   uint
	Questions *[]*QueryQuestion
}

type QueryResult struct {
	Ancount    uint
	Arcount    uint
	Nscount    uint
	RCode      int
	Answers    []*ResourceRecord
	Authority  []*ResourceRecord
	Additional []*ResourceRecord
}

func SearchResourceRecord(query *QueryDomain) (*QueryResult, error) {
	return &QueryResult{}, nil
}

func SearchResourceRecords(query *QueryDomain) (*QueryResult, error) {
	if query.QdCount == 1 {
		return SearchResourceRecord(query)
	}

	return &QueryResult{}, nil
}
