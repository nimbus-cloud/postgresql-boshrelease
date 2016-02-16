package config

import (
	"io/ioutil"
	"github.com/cloudfoundry-incubator/candiedyaml"
)

type Config struct {
	Master     bool            `yaml:"master"`
	Hostname   string          `yaml:"hostname"`
	Ip         string          `yaml:"ip"`
	DnsServers []string        `yaml:"dnsservers"`
	DnsKey     string          `yaml:"dnskey"`
	DnsTTL     int             `yaml:"dnsttl"`
	Logging	   Logging	   `yaml:"logging"`
}

type Logging struct {
	Level 	string 		   `yaml:"level"`
}

func InitConfigFromFile(path string) *Config {

	b, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e.Error())
	}

	c := &Config{}
	e = candiedyaml.Unmarshal(b, c)
	if e != nil {
		panic(e.Error())
	}

	return c
}