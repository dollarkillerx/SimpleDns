package storage

import "github.com/dollarkillerx/easy_dns"

type Interface interface {
	QueryDns(domain string, queryType easy_dns.Type) (resp []easy_dns.Resource, err error)
	StorageDns(domain string, queryType easy_dns.Type, resp *easy_dns.Message, ttl uint32) error
}
