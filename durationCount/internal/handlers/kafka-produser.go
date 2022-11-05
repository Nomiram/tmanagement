package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func KafkaProducer() *kafka.Writer {
	fmt.Println("msg from Producer")
	// to produce messages
	topic := "my-topic-1"
	w := &kafka.Writer{
		Addr: kafka.TCP("kafka:9092"),
		// Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	// if err := w.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }
	return w
}
func KafkaWrite(w *kafka.Writer, key string, value string) {
	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
