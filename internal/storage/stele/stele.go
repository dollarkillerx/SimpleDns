package stele

import (
	"encoding/json"
	"fmt"
	"github.com/dollarkillerx/SimpleDns/pkg/model"
	"github.com/rs/xid"
	"log"
	"strconv"
	"strings"
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
	return []byte(fmt.Sprintf("dns_%s_%d", domain, queryType))
}

func (s *Stele) APIStorageDns(domain string, model *model.DnsModel) error {
	xid := xid.New().String()
	model.ID = xid
	marshal, err := json.Marshal(model)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.stele.Set([]byte(getApiID(xid)), marshal, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	if model.TTL <= 0 {
		model.TTL = 50
	}

	resource, err := iPToAResource(model.Addr)
	if err != nil {
		log.Println(err)
		return err
	}

	msg := &easy_dns.Message{
		Answers: []easy_dns.Resource{
			{
				Header: easy_dns.ResourceHeader{
					Name:   easy_dns.MustNewName(domain + "."),
					Type:   model.QueryType,
					Class:  easy_dns.ClassINET,
					TTL:    uint32(model.TTL),
					Length: 4,
				},
				Body: resource,
			},
		},
	}

	bytes, err := msg.Pack()
	if err != nil {
		log.Println(err)
		return err
	}

	return s.stele.Set(s.getKey(model.Domain+".", model.QueryType), bytes, 0)
}

func (s *Stele) APIDeleteDns(id string) error {
	nodeByte, err := s.stele.Get([]byte(getApiID(id)))
	if err != nil {
		log.Println(err)
		return err
	}

	var mod model.DnsModel
	err = json.Unmarshal(nodeByte, &mod)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.stele.Delete([]byte(getApiID(id)))
	if err != nil {
		log.Println(err)
		return err
	}

	return s.stele.Delete(s.getKey(mod.Domain+".", mod.QueryType))
}

func (s *Stele) APIUpdateDns(id string, model *model.DnsModel) error {
	marshal, err := json.Marshal(model)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.stele.Set([]byte(getApiID(model.ID)), marshal, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	resource, err := iPToAResource(model.Addr)
	if err != nil {
		log.Println(err)
		return err
	}

	if model.TTL <= 0 {
		model.TTL = 50
	}

	msg := &easy_dns.Message{
		Answers: []easy_dns.Resource{
			{
				Header: easy_dns.ResourceHeader{
					Name:   easy_dns.MustNewName(model.Domain + "."),
					Type:   model.QueryType,
					Class:  easy_dns.ClassINET,
					TTL:    uint32(model.TTL),
					Length: 4,
				},
				Body: resource,
			},
		},
	}

	bytes, err := msg.Pack()
	if err != nil {
		log.Println(err)
		return err
	}

	return s.stele.Set(s.getKey(model.Domain+".", model.QueryType), bytes, 0)
}

func (s *Stele) APIListDns() ([]model.DnsModel, error) {
	scan, err := s.stele.PrefixScan([]byte("api_dns."))
	if err != nil {
		return nil, err
	}

	if len(scan) == 0 {
		return []model.DnsModel{}, nil
	}

	var models []model.DnsModel
	for k := range scan {
		var mod model.DnsModel
		err := json.Unmarshal(scan[k].Val, &mod)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		models = append(models, mod)
	}

	return models, nil
}

func getApiID(id string) string {
	return fmt.Sprintf("api_dns.%s", id)
}

func iPToAResource(ip string) (*easy_dns.AResource, error) {
	split := strings.Split(ip, ".")
	if len(split) != 4 {
		return nil, fmt.Errorf("Nof")
	}

	var a [4]byte
	for k, v := range split {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		a[k] = uint8(atoi)
	}

	return &easy_dns.AResource{A: a}, nil
}
