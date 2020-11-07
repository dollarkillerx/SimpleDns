package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Conf struct {
	DNSServer  string `json:"dns_server"`
	ListenAddr string `json:"local_addr"`
}

func ConfInit() *Conf {
	var conf Conf
	file, err := ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.Unmarshal(file, &conf); err != nil {
		log.Fatalln(err)
	}
	return &conf
}
