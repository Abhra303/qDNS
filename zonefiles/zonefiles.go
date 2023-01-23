package zonefiles

import (
	"fmt"

	"github.com/abhra303/qDNS/config"
)

type ResourceRecord interface {
	Match(domain string, flags int) bool
}

type ResourceRecordInfo struct {

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
}

type Soa struct {
	RInfo   ResourceRecordInfo
	MName   string
	RName   string
	Serial  int
	Refresh int
	Retry   int
	Expire  int
	Minimum int
}

type Zonefile struct {
	TTL int
	SOA Soa
	RRs []*ResourceRecord
}

type Zone struct {
	Zonefiles []struct {
		ZonefileBuf *Zonefile
		filePath    string
	}
	ZoneName string
	flags    int
}

type Zones map[string]*Zone

var availableZones Zones = make(Zones, 256)

func loadZones() bool {
	for _, zone := range config.ServerConfiguration.Zones {
		fmt.Printf("%v\n", zone)
	}
	return true
}

func findZone(question *QueryQuestion) (*Zone, bool) {
	if len(availableZones) == 0 && !loadZones() {
		fmt.Printf("zone loading failed\n")
		return &Zone{}, false
	}

	// TODO: fix the searching of available zones
	zone, isPresent := availableZones[question.QName[0]]

	if !isPresent {
		fmt.Println("zone for the query domain not found")
		return &Zone{}, false
	}

	return zone, true
}

func (zone *Zone) findResourceRecord(query *QueryQuestion) QueryResult {
	// code to search for resource record
	return QueryResult{}
}
