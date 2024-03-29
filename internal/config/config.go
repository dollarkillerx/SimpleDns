package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Conf struct {
	DNSServer     string `yaml:"DNSServer"`
	Token         string `yaml:"Token"`
	DNSListenAddr string `yaml:"DNSListenAddr"`
	ApiListenAddr string `yaml:"ApiListenAddr"`
	Debug         bool   `yaml:"Debug"`
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
Token: token
DNSServer: "223.5.5.5:53"
DNSListenAddr: "0.0.0.0:6060"
ApiListenAddr: "0.0.0.0:6061"
Debug: false
`
