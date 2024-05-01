package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

// Config - структура для хранения конфигурации
type Config struct {
	BrokerURL string `json:"brokerURL"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientID  string `json:"clientID"`
	KeepAlive int    `json:"keepAlive"`
	TopicName string `json:"topicName"`
}

// LoadConfig - функция для загрузки конфигурации из файла
func LoadConfig(configFile string) (*Config, error) {
	// Конфигурация по умолчанию
	config := Config{
		BrokerURL: "tcp://localhost:1883",
		Username:  "",
		Password:  "",
		ClientID:  "myClientID",
		KeepAlive: 60, // Значение по умолчанию для keepAlive
		TopicName: "myTopic",
	}

	// Создаем файл конфигурации, если он еще не существует
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Создаем новый файл конфигурации
		configFile, err := os.Create(configFile)
		if err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
		defer configFile.Close()

		// Преобразуем конфигурацию в JSON и записываем в файл
		jsonData, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatalf("Failed to marshal config to JSON: %v", err)
		}
		if _, err := configFile.Write(jsonData); err != nil {
			log.Fatalf("Failed to write config to file: %v", err)
		}
	}

	// Чтение файла конфигурации
	configFileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Декодирование JSON в структуру Config
	err = json.Unmarshal(configFileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return &config, nil
}

func main() {
	// Загрузка конфигурации
	configFile := "config.json"
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	opts := mqtt.NewClientOptions().AddBroker(config.BrokerURL)
	if config.Username != "" && config.Password != "" {
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
	}
	opts.SetClientID(config.ClientID)
	opts.SetKeepAlive(time.Second * time.Duration(config.KeepAlive))
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to MQTT Broker!")
		token := c.Subscribe(config.TopicName, 0, onMessage)
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
