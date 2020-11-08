package routing

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"sync"

	"github.com/dollarkillerx/SimpleDns/internal/storage"
	"github.com/dollarkillerx/easy_dns"
)

// 热更新 路由表
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
	go sto.core()
	return sto
}

func (r *Routing) update() {
	// del old url
	for _, v := range r.old {
		r.storage.DeleteDns(v, easy_dns.TypeA)
	}
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
		r.storage.StorageDns(strings.TrimSpace(split[0]), easy_dns.TypeA, &easy_dns.Message{

		}, 0)
	}
}

func (r *Routing) core() {

}
