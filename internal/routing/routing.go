package routing

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dollarkillerx/SimpleDns/internal/storage"
	"github.com/dollarkillerx/SimpleDns/pkg/utils"
	"github.com/dollarkillerx/easy_dns"
)

// Routing 热更新 路由表
type Routing struct {
	mu sync.Mutex

	oldKey  string
	old     []string
	storage storage.Interface
}

func New(storage storage.Interface) *Routing {
	sto := &Routing{
		old:     []string{},
		storage: storage,
	}

	sto.update()
	go sto.updateRoutingTable()
	return sto
}

func (r *Routing) update() {
	// del old url
	for _, v := range r.old {
		r.storage.DeleteDns(v, easy_dns.TypeA)
	}
	r.old = []string{}
	open, err := os.Open("./routing_table.csv")
	if err != nil {
		os.Create("routing_table.csv")
		return
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		split := strings.Split(string(bytes.TrimSpace(line)), ",")
		if len(split) != 2 {
			continue
		}
		domain := fmt.Sprintf("%s.", strings.TrimSpace(split[0]))

		resource, err := r.iPToAResource(strings.TrimSpace(split[1]))
		if err != nil {
			continue
		}
		r.storage.StorageDns(domain, easy_dns.TypeA, &easy_dns.Message{
			Answers: []easy_dns.Resource{
				{
					Header: easy_dns.ResourceHeader{
						Name:   easy_dns.MustNewName(domain),
						Type:   easy_dns.TypeA,
						Class:  easy_dns.ClassINET,
						TTL:    50,
						Length: 4,
					},
					Body: resource,
				},
			},
		}, 0)
		r.old = append(r.old, domain)
	}
}

func (r *Routing) iPToAResource(ip string) (*easy_dns.AResource, error) {
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

func (r *Routing) updateRoutingTable() {
	for {
		time.Sleep(time.Millisecond * 250)
		file, err := ioutil.ReadFile("./routing_table.csv")
		if err != nil {
			continue
		}
		encode := utils.Md5Encode(string(file))
		if r.oldKey == encode {
			continue
		}
		r.update()
	}
}
