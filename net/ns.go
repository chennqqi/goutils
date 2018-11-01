package net

type NSRecord []string

type NSRecords struct {
	NSRecord
	NicNS map[string]NSRecord `json:"nicns" yaml:"nicns"` //ns set by nic
}
