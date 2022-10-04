package handlers

// Содержит отбработчики для gin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"tmanagement/internal/core"
	"tmanagement/internal/headers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func RedisConnect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
func RedisSet(rdb *redis.Client, key string, value string) {
	var ctx = context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
func RedisGet(rdb *redis.Client, key string) string {
	var ctx = context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

/*
REST API:GET Функция возвращает кратчайшее время для работ

~/duration/<string>
*/
func GetBrowserOptDuration(c *gin.Context) {
	Order_name := c.Param("Order")
	type returnstruct struct {
		Duration float64  `json:"duration"`
		Path     []string `json:"path"`
	}
	//lint:ignore SA4006 (выражение используется далее)
	path := []string{}
	rdb := RedisConnect()
	retstr := RedisGet(rdb, Order_name)
	if retstr != "" {
		var ret returnstruct
		err := json.Unmarshal([]byte(retstr), &ret)
		if err != nil {
			panic(err)
		}
		c.IndentedJSON(http.StatusOK, ret)
		return
	}
	i, path := core.GetOptDuration(Order_name, 10, 200000)
	ret := returnstruct{Duration: i, Path: path}
	if i == -1 {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Status   string
			Duration float64
			Path     []string
		}{fmt.Sprint(http.StatusBadRequest), i, path})
	} else {
		res, _ := json.Marshal(ret)
		RedisSet(rdb, Order_name, string(res))
		c.IndentedJSON(http.StatusOK, ret)
	}
}

// REST API:GET Возвращает информацию из таблицы tasks
func GetTasks(c *gin.Context) {
	Order_name := c.Param("id")

	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from tasks WHERE Order_name = $1", Order_name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tasks := []headers.Task{}

	for rows.Next() {
		p := headers.Task{}
		//Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred
		err := rows.Scan(&p.Task, &p.Order_name, &p.Duration, &p.Resource, &p.Pred)

		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusConflict, err)
			continue
		}
		tasks = append(tasks, p)
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// REST API:GET Возвращает информацию из таблицы orders
func GetOrders(c *gin.Context) {
	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from orders")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	orders := []headers.Order{}

	for rows.Next() {
		p := headers.Order{}
		err := rows.Scan(&p.Order_name, &p.Start_date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, p)
	}
	c.IndentedJSON(http.StatusOK, orders)
}

// -----------FUNCTION postTasks---------------//

// REST API:POST,PUT добавление данных по POST и PUT в таблицу tasks
func PostTasks(c *gin.Context) {
	var newTask headers.Task
	//Получение данных из контекста
	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(newTask)
	//Подключение к Postgres
	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	////Удаление task при обновлении
	if c.Request.Method == "PUT" {
		result, err := db.Exec("DELETE FROM tasks WHERE task = $1 AND order_name = $2; ",
			newTask.Task, newTask.Order_name)
		if err != nil {
			fmt.Println(result)
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
	}
	// Добавление новой работы "task" в таблицу tasks
	result, err := db.Exec("INSERT INTO tasks (task, order_name, duration, resource, pred) values ($1, $2, $3, $4, $5)",
		newTask.Task, newTask.Order_name, newTask.Duration, newTask.Resource, newTask.Pred)
	if err != nil {
		fmt.Println(result)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	// Обнуление данных в redis
	rdb := RedisConnect()
	RedisSet(rdb, newTask.Order_name, "")
	// tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

// ------------FUNCTION postOrders--------------//

// REST API:POST,PUT добавление данных по POST и PUT в таблицу orders
//
// json: {"order_name":string, start_date":string}
func PostOrders(c *gin.Context) {
	var newOrder headers.Order
	//Получение данных
	if err := c.BindJSON(&newOrder); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	////Удаление headers.Order при обновлении
	if c.Request.Method == "PUT" {
		result, err := db.Exec("DELETE FROM orders WHERE order_name = $1; ",
			newOrder.Order_name)
		if err != nil {
			fmt.Println(result)
			c.IndentedJSON(http.StatusBadRequest, err)
			// panic(err)
			return
		}
	}
	//Добавление данных в таблицу
	result, err := db.Exec("insert into orders (order_name, start_date) values ($1, $2)",
		newOrder.Order_name, newOrder.Start_date)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}
	// Обнуление данных в redis
	rdb := RedisConnect()
	RedisSet(rdb, newOrder.Order_name, "")
	c.IndentedJSON(http.StatusCreated, newOrder)

}

// REST API:DELETE удаление данных из таблицы orders
//
// json: {"order_name":string}
func DelOrders(c *gin.Context) {
	var Order headers.Delorder
	//Получение данных
	if err := c.BindJSON(&Order); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Удаление данных из таблицы
	result, err := db.Exec("DELETE FROM orders WHERE order_name = $1",
		Order.Order_name)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}
	// Обнуление данных в redis
	rdb := RedisConnect()
	RedisSet(rdb, Order.Order_name, "")
	c.IndentedJSON(http.StatusCreated, Order)

}

// REST API:DELETE удаление данных из таблицы orders
//
// json: {"order_name":string}
func DelTasks(c *gin.Context) {
	var delTask headers.TaskDel
	//Получение данных
	if err := c.BindJSON(&delTask); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", headers.CONNSTRWDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Удаление данных из таблицы
	result, err := db.Exec("DELETE FROM tasks WHERE order_name = $1 AND task = $2",
		delTask.Order_name, delTask.Task)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}
	// Обнуление данных в redis
	rdb := RedisConnect()
	RedisSet(rdb, delTask.Order_name, "")
	c.IndentedJSON(http.StatusCreated, delTask)

}
