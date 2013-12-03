package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type CheckResult struct {
	Handler     []string `json:"handler"`
	Command     string   `json:"command"`
	Interval    int      `json:"interval"`
	Subscribers []string `json:"subscribers"`
	Standalone  bool     `json:"standalone"`
	Name        string   `json:"name"`
	Issued      int64    `json:"issued"`
	Executed    int64    `json:"executed"`
	Output      string   `json:"output"`
	Status      int      `json:"status"`
	Duration    float32  `json:"duration"`
}

type ResultMsg struct {
	Client string      `json:"client"`
	Check  CheckResult `json:"check"`
}

func runCheck(clientName string, checkName string, checkConf CheckConf, channel chan []byte) {
	// Recover function to log and stop goroutine at panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panicing - "+checkName+" will be ignored until restart: %s\n", err)
			return
		}
	}()

	// Extract binary and args from command
	pathAndArgs := strings.Fields(checkConf.Command)

	for {
		// Create cmd object and extract necessary information
		cmd := exec.Command(pathAndArgs[0], pathAndArgs[1:]...)
		output, _ := cmd.CombinedOutput()
		exitCode := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		// Build structs for marshalling
		checkResult := CheckResult{
			Handler:     checkConf.Handler,
			Command:     checkConf.Command,
			Interval:    checkConf.Interval,
			Subscribers: checkConf.Subscribers,
			Standalone:  checkConf.Standalone,
			Name:        checkName,
			Issued:      time.Now().Unix(),
			Executed:    time.Now().Unix(),
			Output:      string(output),
			Status:      exitCode,
			Duration:    0.0}

		resultMsg := ResultMsg{
			Client: clientName,
			Check:  checkResult}

		// Marshal results to json and send it to channel
		resultJson, err := json.Marshal(resultMsg)
		if err != nil {
			log.Printf("Marshal: %s\n", err)
		}

		channel <- resultJson

		time.Sleep(time.Duration(checkConf.Interval) * time.Second)
	}
}
