package kafka

import (
	"github.com/AzusaChino/ficus/pkg/conf"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var dataCollector sarama.AsyncProducer

func Setup() {
	var err error
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Compression = sarama.CompressionGZIP
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond

	client, err := sarama.NewClient(conf.KafkaConfig.Locations, kafkaConfig)
	if err != nil {
		panic("failed to setup kafka client")
	}
	dataCollector, err = sarama.NewAsyncProducerFromClient(client)
}

func Close() error {
	if err := dataCollector.Close(); err != nil {
		log.Println("Failed to shut down data collector cleanly", err)
	}
	return nil
}

func SendMessage(topic string, key string, value string) {
	dataCollector.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
}