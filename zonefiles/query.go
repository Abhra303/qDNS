package zonefiles

import "fmt"

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

/*
QueryDomain contains the required informations to perform
domain query.
TODO: We should stop storing Questions as pointer
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

func SearchResourceRecord(query *QueryQuestion) (*QueryResult, error) {
	var queryResult QueryResult
	zone, isPresent := findZone(query)
	if !isPresent {
		return nil, fmt.Errorf("zone not found")
	}
	if zone != nil {
		queryResult = zone.findResourceRecord(query)
	}
	return &queryResult, nil
}

func SearchResourceRecords(query *QueryDomain) (*QueryResult, error) {
	if query.QdCount == 1 {
		return SearchResourceRecord((*query.Questions)[0])
	}

	return &QueryResult{}, nil
}
