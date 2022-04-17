package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/jinzhu/configor"
)

var config = struct {
	BotToken           string   `required:"true"`
	ChatID             string   `required:"true"`
	RestartWarpCommand []string `required:"true"`
}{}

func check() bool {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.netflix.com/title/81215567", nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

func getIP() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://ipinfo.io/ip", nil)

	if err != nil {
		return "Error"
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")

	resp, err := client.Do(req)
	if err != nil {
		return "Error"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error"
	}

	return string(body)
}

func sendMsg(text string) {
	client := &http.Client{}

	jsonData := map[string]string{"chat_id": config.ChatID, "text": text}
	jsonBytes, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", "https://api.telegram.org/"+config.BotToken+"/sendMessage", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return
	}
}

func main() {
	err := configor.Load(&config, "config.yml")
	if err != nil {
		log.Panic(err.Error())
	}

	for true {
		if check() {
			time.Sleep(30 * time.Minute)
		} else {
			cmd := exec.Command(config.RestartWarpCommand[0], config.RestartWarpCommand[1:]...)
			cmd.Start()
			cmd.Wait()
			sendMsg("\nNetflix: " + strconv.FormatBool(check()) + "\nNew IP: " + getIP())
		}
	}
}
