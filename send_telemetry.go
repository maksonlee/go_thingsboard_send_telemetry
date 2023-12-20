package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/tidwall/sjson"
)

type Config struct {
	Server   string `yaml:"server"`
	Token    string `yaml:"access_token"`
	Duration int    `yaml:"duration"`
	TS       bool   `yaml:"ts"`
	Message  string `yaml:"message"`
}

func task(c Config) {
	url := c.Server + "/api/v1/" + c.Token + "/telemetry"
	fmt.Println("URL:>", url)
	for range time.Tick(time.Second * 1) {
		var ts = time.Now().UnixMilli()
		var jsonStr []byte
		if c.TS {
			jsonStr = []byte(`{"ts": ` + strconv.FormatInt(ts, 10) + `, "values": ` + c.Message + `}`)
		} else {
			value, _ := sjson.Set(c.Message, "ts", ts)
			jsonStr = []byte(value)
		}
		fmt.Println(bytes.NewBuffer(jsonStr))
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
	}
}

func main() {
	yamlFile, err := os.ReadFile("send_telemetry.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	c := Config{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	go task(c)
	time.Sleep(time.Hour * time.Duration(c.Duration))
}
