package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"tmanagement/internal/handlers"
	"tmanagement/internal/headers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Информация для подключения к БД postgres

func main() {
	value, ok := os.LookupEnv("DBADDR")
	if ok {
		headers.AddDBinCONNSTR(value)
	}

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
	fmt.Println("testing db & redis...")
	if checkDB() {
		fmt.Println("DB ok")

	}
	rdb := handlers.RedisConnect()
	// var ctx = context.Background()
	handlers.RedisSet(rdb, "key", "ok")
	val := handlers.RedisGet(rdb, "key")
	fmt.Println("Redis return " + val)
	router.Run(":8080")
}
func checkDB() bool {
	db, err := sql.Open("postgres", headers.CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	name := "VS"
	res, err := db.Query("SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + name + "'")
	if err != nil {
		fmt.Println(err)
	}
	if res.Next() {
		dbname := ""
		err = res.Scan(&dbname)
		if err != nil {
			panic(err)
		}
		fmt.Println("DB exist:", dbname)
	} else {

		// name := "VS"
		fmt.Println("Creating DB")
		_, err := db.Exec("CREATE DATABASE \"" + name + "\"")
		if err != nil {
			panic(err)
		}
		// fmt.Println(res.RowsAffected())
	}
	db.Close()
	///////////////////
	db, err = sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	b, err := os.ReadFile("init.sql")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	_, err = db.Query(str)
	if err != nil {
		panic(err)
	}
	return true
}
