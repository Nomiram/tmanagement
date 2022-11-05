package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaConsumer() *kafka.Conn {
	fmt.Println("msg from Consumer")

	// to consume messages
	topic := "my-topic-1"
	partition := 0

	// time.Sleep(time.Second * 5)
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		//
		log.Fatal("Comsumer: failed to dial leader:", err)
		time.Sleep(time.Second * 2)
	}

	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// if err := conn.Close(); err != nil {
	// 	log.Fatal("failed to close connection:", err)
	// }
	return conn

}
func KafkaRead(conn *kafka.Conn) (key string, value string) {
	var n kafka.Message
	var err error
	for {
		n, err = conn.ReadMessage(100000)
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		time.Sleep(time.Millisecond * 100)

		if string(n.Key) != "" {
			break
		}
	}
	return string(n.Key), string(n.Value)
}
