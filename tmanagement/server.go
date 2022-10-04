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

/*
TODO:
1. Разбить на файлы
2. Дальнейшая оптимизация кода
*/

// Информация для подключения к БД postgres

// var CONNSTR = "host=db port=5432 user=postgres password=qwerty dbname=VS sslmode=disable"

func main() {
	value, ok := os.LookupEnv("DBADDR")
	if ok {
		headers.AddDBinCONNSTR(value)
		// headers.CONNSTR = "host=" + value + " port=5432 user=postgres password=postgres sslmode=disable"
	}

	// headers.CONNSTR = "host=db port=5432 user=postgres password=postgres sslmode=disable"
	// fmt.Println(headers.CONNSTR)
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
	checkDB()
	fmt.Println("test")
	rdb := handlers.RedisConnect()
	// var ctx = context.Background()
	handlers.RedisSet(rdb, "key", "string")
	val := handlers.RedisGet(rdb, "key")
	fmt.Println(val)
	router.Run(":8080")
}
func checkDB() {
	db, err := sql.Open("postgres", headers.CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	name := "'VS'"
	// name := "\"VS\""
	res, err := db.Query("SELECT datname FROM pg_catalog.pg_database WHERE datname = " + name)
	if err != nil {
		panic(err)
	}
	if !res.Next() {
		name := "VS"
		_, err = db.Exec("CREATE DATABASE " + name)
		if err != nil {
			panic(err)
		}

	}
	db.Close()
	db, err = sql.Open("postgres", headers.CONNSTRWDB)
	// _, err = db.Exec("\\c" + name)
	if err != nil {
		panic(err)
	}
	_, err = db.Query(`
	CREATE TABLE IF NOT EXISTS public.orders
(
    order_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    start_date date,
    CONSTRAINT orders_pkey PRIMARY KEY (order_name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;

-- Table: public.tasks

-- DROP TABLE IF EXISTS public.tasks;

CREATE TABLE IF NOT EXISTS public.tasks
(
    task character varying(10) COLLATE pg_catalog."default" NOT NULL,
    order_name character varying(10) COLLATE pg_catalog."default",
    duration integer,
    resource integer,
    pred character varying(10) COLLATE pg_catalog."default",
    CONSTRAINT tasks_pkey PRIMARY KEY (task, order_name),
    CONSTRAINT tasks_order_name_fkey FOREIGN KEY (order_name)
        REFERENCES public.orders (order_name) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;`)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}
