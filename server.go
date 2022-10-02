package main

import (
	"net/http"
	"tmanagement/internal/handlers"

	// "crypto/rand"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

/*
TODO:
1. Разбить на файлы
2. Дальнейшая оптимизация кода
*/

// Информация для подключения к БД postgres
// var CONNSTR = "user=postgres password=qwerty dbname=VS sslmode=disable"

func main() {
	router := gin.Default()
	router.GET("/duration/:Order", handlers.GetBrowserOptDuration)
	router.GET("/orders", handlers.GetOrders)
	router.GET("/tasks/:id", handlers.GetTasks)
	router.POST("/orders", handlers.PostOrders)
	router.PUT("/orders", handlers.PostOrders)
	router.DELETE("/orders", handlers.DelOrders)
	router.POST("/tasks", handlers.PostTasks)
	router.PUT("/tasks", handlers.PostTasks)
	router.DELETE("/tasks", handlers.DelTasks)
	router.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"ok": "Pong"}) })
	// fmt.Println(getOptDuration("OrderB", 10, 10000))
	router.Run("localhost:8080")
}
