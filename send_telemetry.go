package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	url := `ssl://` + c.Server + `:8883`
	fmt.Println("URL:>", url)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID("")
	opts.SetUsername(c.Token)
	opts.SetPassword("")
	opts.SetCleanSession(false)
	opts.SetTLSConfig(&tls.Config{
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: false})
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Simulator Publisher Started")

	for range time.Tick(time.Second * 1) {
		var ts = time.Now().UnixMilli()
		var jsonStr string
		if c.TS {
			jsonStr = `{"ts": ` + strconv.FormatInt(ts, 10) + `, "values": ` + c.Message + `}`
		} else {
			value, _ := sjson.Set(c.Message, "ts", ts)
			jsonStr = value
		}
		fmt.Println(jsonStr)
		token := client.Publish("v1/devices/me/telemetry", 0, false, jsonStr)
		token.Wait()
	}

	client.Disconnect(250)
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
