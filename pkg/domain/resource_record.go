package domain

type RecordType string

type ResourceRecord struct {
	FQDN  Domain
	ID    string
	Type  RecordType
	Value string
}

const (
	A    RecordType = "A"
	AAAA RecordType = "AAAA"
)
