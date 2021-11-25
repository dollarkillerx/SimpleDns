package main

import (
	"log"

	"github.com/dollarkillerx/SimpleDns/internal/config"
	"github.com/dollarkillerx/SimpleDns/internal/routing"
	"github.com/dollarkillerx/SimpleDns/internal/server"
	"github.com/dollarkillerx/SimpleDns/internal/storage/stele"
)

func main() {
	cfg := config.ConfInit()
	if cfg.Debug {
		log.SetFlags(log.Llongfile | log.LstdFlags)
	}

	stele := stele.New()
	r := routing.New(stele)
	dns := server.New(cfg, stele, r)

	log.Printf("Run In: %s \n", cfg.DNSListenAddr)
	if err := dns.Run(); err != nil {
		log.Fatalln(err)
	}
}
