package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string
	Adress string
	Nick   string
	Pass   string
}

func (config *Config) Parse(configFile string) {
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(*config)
}
