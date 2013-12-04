package main

import (
	"encoding/json"
	"io/ioutil"
)

type CheckConf struct {
	Handler     string    `json:"handler"`
	Command     string    `json:"command"`
	Interval    int       `json:"interval"`
	Standalone  bool      `json:"standalone"`
	Subscribers []string  `json:"subscribers"`
  Occurrences int       `json:"occurrences"`
}

type ClientConf struct {
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	Subscriptions []string  `json:"subscriptions"`
}

type RabbitmqConf struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

func ParseClientConfig(file string) (ClientConf, error) {
	clientJson, err := ioutil.ReadFile(file)
	if err != nil {
		return ClientConf{}, err
	}

	var clientConfig map[string]ClientConf

	if err := json.Unmarshal(clientJson, &clientConfig); err != nil {
		return ClientConf{}, err
	}
	return clientConfig["client"], nil
}

func ParseRabbitmqConfig(file string) (RabbitmqConf, error) {
	rabbitmqJson, err := ioutil.ReadFile(file)
	if err != nil {
		return RabbitmqConf{}, err
	}

	var rabbitmqConfig map[string]RabbitmqConf

	if err := json.Unmarshal(rabbitmqJson, &rabbitmqConfig); err != nil {
		return RabbitmqConf{}, err
	}
	return rabbitmqConfig["rabbitmq"], nil
}

func ParseChecksConfig(file string) (map[string]CheckConf, error) {
	checksJson, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var checksConfig map[string]map[string]CheckConf

	if err := json.Unmarshal(checksJson, &checksConfig); err != nil {
		return nil, err
	}
	return checksConfig["checks"], nil
}
