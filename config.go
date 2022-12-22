package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type User struct {
	Id       int64    `yaml:"id"`
	FrameIds []string `yaml:"frames"`
}

type Journal struct {
	Id    int64  `yaml:"id"`
	Users []User `yaml:"users"`
}

type Config struct {
	Journals []Journal `yaml:"journals"`
}

func loadConfig() {
	configFile, err := os.ReadFile("userFrameConfig.yaml")
	if err != nil {
		log.Fatal("unable to open config; userFrameConfig.yaml")
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("unable to parse config")
	}
}
