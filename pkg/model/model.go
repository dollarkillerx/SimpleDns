package model

import "github.com/dollarkillerx/easy_dns"

type DnsModel struct {
	ID        string        `json:"id"`
	Domain    string        `json:"domain"`
	QueryType easy_dns.Type `json:"query_type"`
	Addr      string        `json:"addr"`
	TTL       int           `json:"ttl"`
}
