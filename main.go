package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
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

func setupInterrupt() {
	// Создание канала для сигналов типа os.Signal. Канал имеет емкость 1.
	c := make(chan os.Signal, 1)

	// Регистрация сигнала SIGINT (обычно генерируется Ctrl+C) для канала c.
	signal.Notify(c, syscall.SIGINT)
	// Начало новой горутины (асинхронного потока выполнения).
	go func() {
		// Блокировка текущей горутины до тех пор, пока не поступит сигнал через канал c.
		<-c
		fmt.Println("Программа прервана пользователем")
		os.Exit(0) // Нормально завершаем программу
	}()
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

	setupInterrupt()
	for {
		time.Sleep(time.Second) // Просто ждем секунду
	}
}
