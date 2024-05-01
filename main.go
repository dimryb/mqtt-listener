package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var clientID string = "myClientID"
var brokerURL string = "tcp://localhost:1883" // Замените на IP адрес или домен вашего брокера MQTT
var username string = ""
var password string = ""

const topicName = "myTopic"

func onMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message arrived: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	opts.SetClientID(clientID)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to MQTT Broker!")
		token := c.Subscribe(topicName, 0, onMessage)
		token.Wait()
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("Connection lost!")
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	<-time.After(1 * time.Second) // Ожидание для демонстрации получения сообщений
}
