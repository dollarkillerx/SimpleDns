package stele

import (
	"fmt"
	"log"
	"time"

	"github.com/dollarkillerx/SimpleDns/internal/storage"
	"github.com/dollarkillerx/easy_dns"
	"github.com/dollarkillerx/stele/pkg/stele"
)

type Stele struct {
	stele *stele.Local
}

func New() storage.Interface {
	stele, err := stele.NewLocal("./stele_data")
	if err != nil {
		log.Fatalln(err)
	}
	output := &Stele{
		stele: stele,
	}

	return output
}

func (s *Stele) QueryDns(domain string, queryType easy_dns.Type) (resp []easy_dns.Resource, err error) {
	var m easy_dns.Message
	data, err := s.stele.Get(s.getKey(domain, queryType))
	if err != nil {
		return nil, err
	}

	if err := m.Unpack(data); err != nil {
		return nil, err
	}

	return m.Answers, nil
}

func (s *Stele) StorageDns(domain string, queryType easy_dns.Type, resp *easy_dns.Message, ttl uint32) error {
	pack, err := resp.Pack()
	if err != nil {
		return err
	}
	var tt int64
	if ttl != 0 {
		tt = int64(ttl) * int64(time.Second)
	}
	return s.stele.Set(s.getKey(domain, queryType), pack, tt)
}

func (s *Stele) DeleteDns(domain string, queryType easy_dns.Type) error {
	return s.stele.Delete(s.getKey(domain, queryType))
}

func (s *Stele) getKey(domain string, queryType easy_dns.Type) []byte {
	return []byte(fmt.Sprintf("%s_%d", domain, queryType))
}
