package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/3011/project/findCloudflareIP/scanner/config"
	"github.com/3011/project/findCloudflareIP/scanner/db"
	"github.com/bitly/go-simplejson"
)

type Prefixe struct {
	Netblock string
	Size     int
}

type Colo struct {
	City      string
	Region    string
	Continent string
}

var (
	client  = &http.Client{Timeout: 3 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}
	coloMap = make(map[string]Colo)
)

func Dec2IP(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func IP2Dec(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func reqGet(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	req.Header.Set("Accept", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func isCloudflare(ip string) {
	body, err := reqGet("http://" + ip)
	if err != nil {
		return
	}
	if !strings.Contains(body, "Cloudflare Ray ID") || !strings.Contains(body, "Direct IP access not allowed") {
		return
	}

	body, err = reqGet("http://" + ip + "/cdn-cgi/trace")
	if err != nil {
		return
	}
	if !strings.Contains(body, "warp=off") {
		return
	}
	bodyLines := strings.Split(body, "\n")
	for _, line := range bodyLines {
		if strings.HasPrefix(line, "colo=") {
			colo := strings.Replace(line, "colo=", "", 1)
			saveIP(ip, colo)
			break
		}
	}
}

func saveIP(ip string, colo string) {
	info := coloMap[colo]
	sendMessage(ip + " " + info.City + ", " + info.Region)
	db.CreateEndPoint(&db.EndPoint{IP: ip, Colo: colo, City: info.City, Region: info.Region, Continent: info.Continent})
}

func sendMessage(text string) {
	url := "https://api.telegram.org/" + config.Config.BotToken + "/sendMessage"
	jsonStr := []byte(`{"chat_id":"` + config.Config.ChatID + `","text":"` + text + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		sendMessage("sendmsg fail")
	}
	defer resp.Body.Close()
}

func worker(ips <-chan string) {
	for ip := range ips {
		isCloudflare(ip)
	}
}

func readJSON() []Prefixe {
	file, err := os.Open("as.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	asnJSON, err := simplejson.NewJson(content)
	if err != nil {
		panic(err)
	}
	arr, err := asnJSON.Get("prefixes").Array()
	temp := make(map[string]int)
	for _, v := range arr {
		prefixe, ok := v.(map[string]interface{})
		if !ok {
			panic("json err")
		}
		netblock := strings.Split(prefixe["netblock"].(string), "/")[0]
		size, _ := strconv.Atoi(prefixe["size"].(string))

		if _, ok := temp[netblock]; ok {
			if size > temp[netblock] {
				temp[netblock] = size
			}
		} else {
			temp[netblock] = size
		}
	}
	var prefixes []Prefixe
	for n, s := range temp {
		prefixes = append(prefixes, Prefixe{Netblock: n, Size: s})
	}
	return prefixes
}

func init() {
	file, err := os.Open("colo.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	coloJSON, err := simplejson.NewJson(content)
	if err != nil {
		panic(err)
	}
	arr, err := coloJSON.Map()
	if err != nil {
		panic(err)
	}

	for k, v := range arr {
		colo := v.(map[string]interface{})
		coloMap[k] = Colo{City: colo["city"].(string), Region: colo["region"].(string), Continent: colo["continent"].(string)}
	}
}

func main() {
	workerNum := 256
	ips := make(chan string, 2)

	for i := 0; i < workerNum; i++ {
		go worker(ips)
	}

	prefixes := readJSON()
	for _, v := range prefixes {
		decIP := IP2Dec(v.Netblock)
		for i := 0; i <= v.Size; i++ {
			ip := Dec2IP(decIP + int64(i))
			ips <- ip
		}
	}

	time.Sleep(10 * time.Second)
	sendMessage("Done!")
}
