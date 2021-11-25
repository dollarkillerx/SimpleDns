package server

import (
	"log"
	"net"
	"time"

	"github.com/dollarkillerx/SimpleDns/internal/config"
	"github.com/dollarkillerx/SimpleDns/internal/routing"
	"github.com/dollarkillerx/SimpleDns/internal/storage"
	"github.com/dollarkillerx/easy_dns"
	"github.com/gin-gonic/gin"
)

type SimpleDns struct {
	conf    *config.Conf
	storage storage.Interface
	r       *routing.Routing
	app     *gin.Engine
}

func New(conf *config.Conf, storage storage.Interface, r *routing.Routing) *SimpleDns {
	app := gin.New()
	app.Use(gin.Recovery())
	if conf.Debug {
		app.Use(gin.Logger())
	}

	sim := &SimpleDns{
		conf:    conf,
		storage: storage,
		r:       r,
		app:     app,
	}

	return sim
}

func (s *SimpleDns) Run() error {
	s.api()

	addr, err := net.ResolveUDPAddr("udp", s.conf.DNSListenAddr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		i, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go s.core(buf[:i], addr, conn)
	}
}

func (s *SimpleDns) core(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	var msg easy_dns.Message
	if err := msg.Unpack(data); err != nil {
		log.Println(err)
		return
	}
	if len(msg.Questions) <= 0 {
		return
	}
	if len(msg.Questions) > 1 {
		// 直接拨号 通常情况都是小于1的
		dns, err := s.dialDns(data)
		if err != nil {
			log.Println(err)
			return
		}
		pack, err := dns.Pack()
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := conn.WriteToUDP(pack, addr); err != nil {
			log.Println(err)
			return
		}
		return
	}

	// 检测本地路由表
	// 检测缓存
	dns, err := s.storage.QueryDns(msg.Questions[0].Name.String(), msg.Questions[0].Type)
	if err == nil {
		msg.Header.Response = true
		msg.Answers = dns

		pack, err := msg.Pack()
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := conn.WriteToUDP(pack, addr); err != nil {
			log.Println(err)
			return
		}

		return
	}

	// 发起DNS拨号
	dnsResp, err := s.dialDns(data)
	if err != nil {
		log.Println(err)
		return
	}
	pack, err := dnsResp.Pack()
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(pack, addr); err != nil {
		log.Println(err)
		return
	}

	// storage
	if err := s.storage.StorageDns(msg.Questions[0].Name.String(), msg.Questions[0].Type, dnsResp, dnsResp.Answers[0].Header.TTL); err != nil {
		log.Println(err)
	}
}

func (s *SimpleDns) dialDns(msg []byte) (*easy_dns.Message, error) {
	dial, err := net.Dial("udp", s.conf.DNSServer)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dial.SetWriteDeadline(time.Now().Add(time.Second))
	dial.SetReadDeadline(time.Now().Add(time.Second))
	dial.SetDeadline(time.Now().Add(time.Second))
	_, err = dial.Write(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	buf := make([]byte, 512)
	var m easy_dns.Message
	read, err := dial.Read(buf)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := m.Unpack(buf[:read]); err != nil {
		log.Println(err)
		return nil, err
	}

	return &m, nil
}
