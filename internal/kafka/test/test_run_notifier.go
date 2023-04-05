package main

import (
	"kurushimi/internal/config"
	"kurushimi/internal/kafka"
)

func main() {
	writer := kafka.NewKafkaNotifier(&config.KafkaConfig{
		Host: "localhost",
		Port: 9092,
	})

	//err := writer(context.Background())
	//if err != nil {
	//	panic(err)
	//}
}
