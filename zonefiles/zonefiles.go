package zonefiles

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/abhra303/qDNS/config"
	"github.com/abhra303/qDNS/ds/trie"
)

type RType uint16
type RClass uint16

const (
	UnknownType RType = iota
	SOA
	NS
	A
	Aaaa
	MX
	TXT
	Cname
)

const (
	UnknownClass RClass = iota
	IN
	CS
	HS
)

var Catalog trie.Trie

type ResourceRecord interface {
	GetRType() RType
	GetRClass() RClass
	GetValue() string
	GetTtl() uint
}

type resourceRecord struct {

	/*
	   Two octets containing one of the RR type codes.  This
	   field specifies the meaning of the data in the RDATA
	   field.
	*/
	Type RType

	/*
	   two octets which specify the class of the data in the
	   RDATA field.
	*/
	Class RClass

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
	Value string
}

/*
Though these records have similar data, we

	should still store them in different struct.
	This will help us if we have to do any record-
	specific changes.
*/
type ARecord struct {
	resourceRecord
}

func (a *ARecord) GetRClass() RClass {
	return a.Class
}

func (a *ARecord) GetRType() RType {
	return A
}

func (a *ARecord) GetValue() string {
	return a.Value
}

func (a *ARecord) GetTtl() uint {
	return a.TTL
}

type AaaaRecord struct {
	resourceRecord
}

func (aaaa *AaaaRecord) GetRClass() RClass {
	return aaaa.Class
}

func (aaaa *AaaaRecord) GetRType() RType {
	return Aaaa
}

func (aaaa *AaaaRecord) GetValue() string {
	return aaaa.Value
}

func (aaaa *AaaaRecord) GetTtl() uint {
	return aaaa.TTL
}

type NSRecord struct {
	resourceRecord
}

func (n *NSRecord) GetRClass() RClass {
	return n.Class
}

func (n *NSRecord) GetRType() RType {
	return NS
}

func (n *NSRecord) GetValue() string {
	return n.Value
}

func (n *NSRecord) GetTtl() uint {
	return n.TTL
}

type TxtRecord struct {
	resourceRecord
}

func (t *TxtRecord) GetRClass() RClass {
	return t.Class
}

func (t *TxtRecord) GetRType() RType {
	return TXT
}

func (t *TxtRecord) GetValue() string {
	return t.Value
}

func (t *TxtRecord) GetTtl() uint {
	return t.TTL
}

type CnameRecord struct {
	resourceRecord
}

func (c *CnameRecord) GetRClass() RClass {
	return c.Class
}

func (c *CnameRecord) GetRType() RType {
	return Cname
}

func (c *CnameRecord) GetValue() string {
	return c.Value
}

func (c *CnameRecord) GetTtl() uint {
	return c.TTL
}

type MxRecord struct {
	resourceRecord
	Preference int
}

func (m *MxRecord) GetRClass() RClass {
	return m.Class
}

func (m *MxRecord) GetRType() RType {
	return MX
}

func (m *MxRecord) GetValue() string {
	return m.Value
}

func (m *MxRecord) GetTtl() uint {
	return m.TTL
}

func (m *MxRecord) GetPreference() int {
	return m.Preference
}

type Soa struct {
	Class   RClass
	MName   string
	RName   string
	Serial  int
	Refresh int
	Retry   int
	Expire  int
	Minimum int
}

type Zone struct {
	trie     trie.Trie
	ZoneName string
	TTL      int
	SOA      Soa
	Origin   string
	flags    int32
}

func (z *Zone) Put(key string, data interface{}) error {
	return nil
}

func (z *Zone) Update(key string, data interface{}) error {
	return nil
}

func (z *Zone) Delete(key string) (interface{}, error) {
	return nil, nil
}

func (z *Zone) Search(key string) (interface{}, error) {
	return nil, nil
}

func (z *Zone) IsEmpty() bool {
	return z.trie.IsEmpty()
}

type zonefileParser struct {
	zone          *Zone
	fscanner      *bufio.Scanner
	currentDomain string
}

var domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func CheckDomainValidity(str string) bool {
	return domainRegexp.MatchString(str)
}

func CheckIPv4Validity(str string) bool {
	if ip := net.ParseIP(str); ip != nil {
		return strings.Count(str, ".") == 4 && strings.Count(str, ":") == 0
	}
	return false
}

func CheckIPv6Validity(str string) bool {
	if ip := net.ParseIP(str); ip != nil {
		return strings.Count(str, ":") > 2
	}
	return false
}

func CheckClassValidity(str string) RClass {
	switch str {
	case "IN":
		return IN
	case "CS":
		return CS
	case "HS":
		return HS
	}
	return UnknownClass
}

func (zp *zonefileParser) parseZoneDirectives(line string) error {
	if strings.HasPrefix(line, "$") {
		data := strings.Split(line, "=")
		if len(data) > 2 {
			return fmt.Errorf("invalid file: can't have more than one equals for directives \"%v\"", line)
		} else if len(data) == 1 {
			return fmt.Errorf("invalid file: \"=\" missing in directive declaration \"%v\"", line)
		}
		rawValue := strings.Split(data[1], ";")
		value := strings.TrimSpace(rawValue[0])
		data[0] = strings.TrimLeft(data[0], "$")
		switch data[0] {
		case "ORGIN":
			if CheckDomainValidity(value) {
				return fmt.Errorf("invalid file: the value for ORIGIN is not valid \"%v\"", line)
			}
			if !strings.HasSuffix(value, ".") {
				value += "."
			}
			zp.zone.Origin = value
		case "TTL":
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil
			}
			zp.zone.TTL = i
		default:
			return fmt.Errorf("invalid file: unknown directive \"%v\"", line)
		}
	} else {
		return fmt.Errorf("invalid file: the line doesn't contain directive \"%v\"", line)
	}
	return nil
}

func (zp *zonefileParser) parseMetadataFromLine(fields []string, resoresourceRecord *resourceRecord) error {
	fieldNumbers := len(fields)
	var class RClass = IN

	// for i := 0; i < fieldNumbers-1 && checkTypeValidity(fields[i]) == UnknownType; i++ {
	// 	switch fields[i] {
	// 	case "@":
	// 		zp.currentDomain = zp.zone.Origin
	// 	case
	// 	}
	// }

	if fieldNumbers == 2 {
		if zp.currentDomain == "" {
			return fmt.Errorf("invalid file: missing domain name")
		}
	} else if fieldNumbers == 3 {
		if fields[0] == "@" {
			zp.currentDomain = zp.zone.Origin
		} else if CheckDomainValidity(fields[0]) {
			zp.currentDomain = fields[0]
		} else if class = CheckClassValidity(fields[0]); class != UnknownClass {
		} else {
			return fmt.Errorf("invalid file: unknown first field")
		}
	} else if fieldNumbers == 4 {
		if fields[0] == "@" {
			zp.currentDomain = zp.zone.Origin
		} else if !CheckDomainValidity(fields[0]) {
			return fmt.Errorf("invalid file: the domain name is not valid")
		}
		zp.currentDomain = fields[0]
		class = CheckClassValidity(fields[1])
		if class == UnknownClass {
			return fmt.Errorf("invalid file: unknown class field")
		}
	}
	resoresourceRecord.Class = class

	return nil
}

func (zp *zonefileParser) parseSoaFromFile(fields []string) error {
	return nil
}

func (zp *zonefileParser) parseNsFromFile(fields []string) error {
	var err error
	var value string
	fieldNumbers := len(fields)
	nsRecord := NSRecord{resourceRecord: resourceRecord{Type: NS, TTL: uint(zp.zone.TTL)}}

	err = zp.parseMetadataFromLine(fields, &nsRecord.resourceRecord)
	if err != nil {
		return err
	}
	if !CheckDomainValidity(fields[fieldNumbers-1]) {
		return fmt.Errorf("invalid file: the cname value is not valid")
	}
	value = fields[fieldNumbers-1]
	nsRecord.Value = value
	zp.zone.Put(zp.currentDomain, nsRecord)
	return nil
}

func (zp *zonefileParser) parseAFromFile(fields []string) error {
	fieldNumbers := len(fields)
	var err error
	var value string
	aRecord := ARecord{resourceRecord: resourceRecord{Type: A}}

	err = zp.parseMetadataFromLine(fields, &aRecord.resourceRecord)
	if err != nil {
		return err
	}
	if !CheckIPv4Validity(fields[fieldNumbers-1]) {
		return fmt.Errorf("invalid file: the given ipv4 is not valid")
	}
	value = fields[fieldNumbers-1]
	aRecord.Value = value
	zp.zone.Put(zp.currentDomain, aRecord)
	return nil
}

func (zp *zonefileParser) parseAaaaFromFile(fields []string) error {
	fieldNumbers := len(fields)
	var err error
	var value string
	aaaaRecord := AaaaRecord{resourceRecord: resourceRecord{Type: Aaaa}}

	err = zp.parseMetadataFromLine(fields, &aaaaRecord.resourceRecord)
	if err != nil {
		return err
	}
	if !CheckIPv6Validity(fields[fieldNumbers-1]) {
		return fmt.Errorf("invalid file: the given ipv6 is not valid")
	}
	value = fields[fieldNumbers-1]
	aaaaRecord.Value = value
	zp.zone.Put(zp.currentDomain, aaaaRecord)
	return nil
}

func (zp *zonefileParser) parseMxFromFile(fields []string) error {
	fieldNumbers := len(fields)
	var class RClass = IN
	var preference int

	if fieldNumbers == 2 {
		if zp.currentDomain == "" {
			return fmt.Errorf("invalid file: missing domain name")
		}
	} else if fieldNumbers == 3 {
		if fields[0] == "@" {
			zp.currentDomain = zp.zone.Origin
		} else if CheckDomainValidity(fields[0]) {
			zp.currentDomain = fields[0]
		} else if class = CheckClassValidity(fields[0]); class != UnknownClass {
		} else if i, err := strconv.Atoi(fields[1]); err == nil {
			preference = i
		} else {
			return fmt.Errorf("invalid file: the mx record is invalid")
		}
	} else if fieldNumbers == 4 {
		if i, err := strconv.Atoi(fields[2]); err != nil {
			if fields[0] == "@" {
				zp.currentDomain = zp.zone.Origin
			} else if !CheckDomainValidity(fields[0]) {
				return fmt.Errorf("invalid file: the domain name is not valid")
			}
			zp.currentDomain = fields[0]
			class := CheckClassValidity(fields[1])
			if class == UnknownClass {
				return fmt.Errorf("invalid file: unknown class field")
			}
		} else {
			preference = i
			if fields[0] == "@" {
				zp.currentDomain = zp.zone.Origin
			} else if CheckDomainValidity(fields[0]) {
				zp.currentDomain = fields[0]
			} else if class = CheckClassValidity(fields[0]); class != UnknownClass {
			} else {
				return fmt.Errorf("invalid file: the mx record is invalid")
			}
		}
	} else if fieldNumbers == 5 {
		if fields[0] == "@" {
			zp.currentDomain = zp.zone.Origin
		} else if !CheckDomainValidity(fields[0]) {
			return fmt.Errorf("invalid file: the domain name is not valid")
		}
		zp.currentDomain = fields[0]
		class = CheckClassValidity(fields[1])
		if class == UnknownClass {
			return fmt.Errorf("invalid file: unknown class field")
		}
		i, err := strconv.Atoi(fields[3])
		if err != nil {
			return err
		}
		preference = i
	}
	if !CheckDomainValidity(fields[fieldNumbers-1]) {
		return fmt.Errorf("invalid file: the value of mx field is invalid")
	}
	value := fields[fieldNumbers-1]

	mxRecord := MxRecord{resourceRecord: resourceRecord{Type: MX, Class: class, Value: value}, Preference: preference}
	zp.zone.Put(zp.currentDomain, mxRecord)
	return nil
}

func (zp *zonefileParser) parseTxtFromFile(fields []string) error {
	var err error
	var value string
	txtRecord := TxtRecord{resourceRecord: resourceRecord{Type: TXT}}

	err = zp.parseMetadataFromLine(fields, &txtRecord.resourceRecord)
	if err != nil {
		return err
	}
	value = fields[len(fields)-1]
	txtRecord.Value = value
	zp.zone.Put(zp.currentDomain, txtRecord)
	return nil
}

func (zp *zonefileParser) parseCnameFromFile(fields []string) error {
	var err error
	var value string
	fieldNumbers := len(fields)
	cnameRecord := CnameRecord{resourceRecord: resourceRecord{Type: Cname}}

	err = zp.parseMetadataFromLine(fields, &cnameRecord.resourceRecord)
	if err != nil {
		return err
	}
	if !CheckDomainValidity(fields[fieldNumbers-1]) {
		return fmt.Errorf("invalid file: the cname value is not valid")
	}
	value = fields[fieldNumbers-1]
	cnameRecord.Value = value
	zp.zone.Put(zp.currentDomain, cnameRecord)
	return nil
}

func (zp *zonefileParser) getRrType(fields []string) RType {
	for i := 0; i < len(fields)-1; i++ {
		if fields[i] == "A" {
			return A
		} else if fields[i] == "AAAA" {
			return Aaaa
		} else if fields[i] == "NS" {
			return NS
		} else if fields[i] == "MX" {
			return MX
		} else if fields[i] == "TXT" {
			return TXT
		} else if fields[i] == "CNAME" {
			return Cname
		} else if fields[i] == "SOA" {
			return SOA
		}
	}
	return UnknownType
}

func (zp *zonefileParser) parseFile() error {
	var err error
	fscanner := zp.fscanner

	for fscanner.Scan() {
		line := fscanner.Text()

		if strings.HasPrefix(line, ";") || line == "" {
			continue
		}

		if strings.HasPrefix(line, "$") {
			err = zp.parseZoneDirectives(line)
			if err != nil {
				return err
			}
			continue
		}

		line = strings.Split(line, ";")[0]
		fields := strings.Fields(line)
		fieldNumbers := len(fields)
		if fieldNumbers == 0 {
			continue
		} else if fieldNumbers == 1 {
			return fmt.Errorf("invalid file: record has only one field \"%v\"", line)
		}

		switch zp.getRrType(fields) {
		case SOA:
			err = zp.parseSoaFromFile(fields)
		case NS:
			err = zp.parseNsFromFile(fields)
		case A:
			err = zp.parseAFromFile(fields)
		case Aaaa:
			err = zp.parseAaaaFromFile(fields)
		case MX:
			err = zp.parseMxFromFile(fields)
		case TXT:
			err = zp.parseTxtFromFile(fields)
		case Cname:
			err = zp.parseCnameFromFile(fields)
		case UnknownType:
			return fmt.Errorf("unable to parse resource type")
		}
		if err != nil {
			return err
		}
	}
	if err = fscanner.Err(); err != nil {
		return err
	}
	return nil
}

func (z *Zone) findResourceRecord(query *QueryQuestion) QueryResult {
	// code to search for resource record
	return QueryResult{}
}

func (z *Zone) loadFromFiles(files []string) error {
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		fscanner := bufio.NewScanner(f)
		zfParser := zonefileParser{zone: z, fscanner: fscanner}

		err = zfParser.parseFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func loadZones() bool {
	for _, zoneConf := range config.ServerConfiguration.Zones {
		var zone *Zone
		fmt.Printf("%v\n", zoneConf)
		zone.trie = trie.NewTrie(&trie.TrieContext{})
		zone.ZoneName = zoneConf.ZoneName
		err := zone.loadFromFiles(zoneConf.ZonefileLocation)
		if err != nil {
			fmt.Println(fmt.Errorf("error loading zones: %v", err))
			return false
		}
		Catalog.Put(zone.ZoneName, zone)
	}
	return true
}

func findZone(question *QueryQuestion) (*Zone, bool) {
	if Catalog.IsEmpty() && !loadZones() {
		return nil, false
	}

	// TODO: fix the searching of available zones
	zone, err := Catalog.Search(question.QName)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	fmt.Println("zone for the query domain not found")
	return zone.(*Zone), true
}
