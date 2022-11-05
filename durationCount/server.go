package main

import (
	"durationCount/internal/core"
	"durationCount/internal/handlers"
	"durationCount/internal/headers"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	// "github.com/confluentinc/confluent-kafka-go"
)

func main() {
	value, ok := os.LookupEnv("DBADDR")
	if ok {
		headers.AddDBinCONNSTR(value)
	}
	cons := handlers.KafkaConsumer()
	writer := handlers.KafkaProducer()
	router := gin.Default()
	router.GET("/duration/:Order", handlers.GetBrowserOptDuration)
	router.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"ok": "Pong"}) })
	// fmt.Println(getOptDuration("OrderB", 10, 10000))
	// fmt.Println("test")
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("read")
		msg_type, order_name := handlers.KafkaRead(cons)
		if msg_type[:5] != "input" {
			continue
		}
		fmt.Println("order_name: ", order_name)
		type returnstruct struct {
			Duration float64
			Path     []string
		}
		mintime, minpath := core.GetOptDuration(order_name, 10, 100000)
		str, _ := json.Marshal(struct {
			int
			returnstruct
		}{200, returnstruct{mintime, minpath}})
		handlers.KafkaWrite(writer, "return"+msg_type[5:], string(str))
		// router.Run(":6000")
	}
}
