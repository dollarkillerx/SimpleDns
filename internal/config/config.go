package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Conf struct {
	DNSServer  string `yaml:"DNSServer"`
	ListenAddr string `yaml:"ListenAddr"`
	Debug      bool   `yaml:"Debug"`
}

func ConfInit() *Conf {
	var conf Conf
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		create, err := os.Create("./config.yaml")
		if err != nil {
			log.Fatalln(err)
		}
		defer create.Close()
		create.WriteString(base)
		log.Fatalln("NOT CONFIG")
	}

	if err := yaml.Unmarshal(file, &conf); err != nil {
		log.Fatalln(err)
	}
	return &conf
}

var base = `
DNSServer: "223.5.5.5:53"
ListenAddr: "0.0.0.0:53"
Debug: false`
