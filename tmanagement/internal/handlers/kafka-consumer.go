package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaConsumer() *kafka.Conn {
	fmt.Println("msg from Consumer")

	// to consume messages
	topic := "my-topic-1"
	partition := 0
	var conn *kafka.Conn
	var err error
	// time.Sleep(time.Second * 5)
	for {
		conn, err = kafka.DialLeader(context.Background(), "tcp", "kafka_0:9092", topic, partition)
		if err != nil {
			//
			fmt.Println("Comsumer: failed to dial leader:", err)
			time.Sleep(time.Second * 2)
			continue
		}
		break
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
