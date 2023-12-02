package main

import (
	"os"

	"github.com/bp1222/photoSync/mail"
	"github.com/bp1222/photoSync/tinybeans"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Sender struct {
	From  string           `yaml:"from"`
	Gmail *mail.GmailEmail `yaml:"gmail"`
	Smtp  *mail.SMTPEmail  `yaml:"smtp"`
}

type Mitm struct {
	Host string `yaml:"host"`
	Port *int   `yaml:"port"`
}

type Config struct {
	Live      bool             `yaml:"live"`
	Mitm      *Mitm            `yaml:"mitm"`
	Sender    Sender           `yaml:"sender"`
	Tinybeans tinybeans.Config `yaml:"tinybeans"`
}

func loadConfig() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("unable to open config; config.yaml")
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("unable to parse config")
	}
}
