package storage

import (
	"github.com/dollarkillerx/SimpleDns/pkg/model"
	"github.com/dollarkillerx/easy_dns"
)

type Interface interface {
	QueryDns(domain string, queryType easy_dns.Type) (resp []easy_dns.Resource, err error)
	DeleteDns(domain string, queryType easy_dns.Type) error
	StorageDns(domain string, queryType easy_dns.Type, resp *easy_dns.Message, ttl uint32) error

	APIStorageDns(domain string, model *model.DnsModel) error
	APIDeleteDns(id string) error
	APIUpdateDns(id string, model *model.DnsModel) error
	APIListDns() ([]model.DnsModel, error)
}
