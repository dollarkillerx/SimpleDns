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
	dns := server.New(cfg, stele)
	routing.New(stele)

	log.Printf("Run In: %s \n", cfg.ListenAddr)
	if err := dns.Run(); err != nil {
		log.Fatalln(err)
	}
}
