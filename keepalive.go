package main

import (
	"encoding/json"
	"time"
)

// Struct with necessary information for keepalive event
type KeepAliveMsg struct {
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	Subscriptions []string `json:"subscriptions"`
	Timestamp     int64    `json:"timestamp"`
}

func (k *KeepAliveMsg) setTime() {
	k.Timestamp = time.Now().Unix()
}

func KeepAlive(clientConf ClientConf, goChannel chan []byte) {
	interval := 60
	msg := &KeepAliveMsg{clientConf.Name, clientConf.Address, clientConf.Subscriptions, 0}
	for {
		msg.setTime()
		msgJson, _ := json.Marshal(msg)
		goChannel <- msgJson
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
