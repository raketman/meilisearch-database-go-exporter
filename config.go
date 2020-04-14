package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
)
type Work struct {
	Index string `json:"index""`
	Primary string `json:"primary"`
	DeleteBefore bool `json:"delete_before"`
	DB_DRIVER string `json:"db_driver"`
	DB_DSN string `json:"db_dsn"`
	Query string `json:"query"`
	Sleep int `json:"sleep"`
	Thread int `json:"thread"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	DisplayedAttributes []string `json:"displayed_attributes"`
	SearchableAttributes []string `json:"searchable_attributes"`
}

type Config struct {
	Host string  `json:"host""`
	Key string `json:"key""`
	Works []Work  `json:"works"`



}

func (c *Config) Read(configFile string) {

	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("ReadConfig", err)
	}

	err = json.Unmarshal(buf, c)
	if err != nil {
		log.Fatal("UnmarshalConfig:", err)
	}
}
