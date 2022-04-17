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

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func check() bool {
	for {
		req, err := http.NewRequest("GET", "https://www.netflix.com/title/81215567", nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		return resp.StatusCode == 200
	}
}

func getIP() string {
	for {
		req, err := http.NewRequest("GET", "http://ipinfo.io/ip", nil)

		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		return string(body)
	}
}

func sendMsg(text string) {

	jsonData := map[string]string{"chat_id": config.ChatID, "text": text}
	jsonBytes, _ := json.Marshal(jsonData)

	for {
		req, err := http.NewRequest("POST", "https://api.telegram.org/"+config.BotToken+"/sendMessage", bytes.NewBuffer(jsonBytes))
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
		req.Header.Set("Content-Type", "application/json")

		_, err = client.Do(req)
		if err != nil {
			continue
		}
	}
}

func main() {
	err := configor.Load(&config, "config.yml")
	if err != nil {
		log.Panic(err.Error())
	}
	sendMsg("Netflix: " + strconv.FormatBool(check()) + "\nNew IP: " + getIP())

	for {
		log.Println("for")
		if check() {
			log.Println("yes")
			time.Sleep(30 * time.Minute)
		} else {
			log.Println("no")
			cmd := exec.Command(config.RestartWarpCommand[0], config.RestartWarpCommand[1:]...)
			log.Println(1)
			cmd.Run()

			log.Println(2)
			sendMsg("Netflix: " + strconv.FormatBool(check()) + "\nNew IP: " + getIP())
			log.Println("end")
		}
	}
}
