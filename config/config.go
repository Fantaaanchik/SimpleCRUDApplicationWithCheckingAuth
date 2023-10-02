package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	App              App    `json:"app"`
	TokenLifeSeconds int    `json:"token_life_seconds"`
	Secret           string `json:"secret"`
	DB               string `json:"db"`
}

type App struct {
	ServerName string `json:"server_name"`
	PortRun    string `json:"port_run"`
}

var Conf Config

func ReadConf(F string) {
	byteValue, err := ioutil.ReadFile(F)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	err = json.Unmarshal(byteValue, &Conf)
	//fmt.Println(Config)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

}
