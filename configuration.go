package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type AMQPConfiguration struct {
	Uri                string `json:"uri"`
}

type DataBaseConfiguration struct {
	Dialect            string `json:"dialect"`
	Uri                string `json:"uri"`
	MaxIdleConnections int    `json:"max-idle-connections"`
	MaxOpenConnections int    `json:"max-open-connections"`
	ShowLog            bool   `json:"log"`
	Threads            uint8  `json:"threads"`
	Limit              uint16 `json:"limit"`
}

type Configuration struct {
	DB   DataBaseConfiguration `json:"db"`
	AMQP AMQPConfiguration     `json:"amqp"`
}

func (this *Configuration) Init(configFile string) {
	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(configJson, &this)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
