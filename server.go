package main

import (
	// "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var CONNSTR = "user=postgres password=qwerty dbname=VS sslmode=disable"

type order struct {
	Order_name string `json:"order_name"`
	Start_date string `json:"start_date"`
}
type task struct {
	Task       string `json:"task"`
	Order_name string `json:"order_name"`
	Duration   int    `json:"duration"`
	Resource   int    `json:"resource"`
	Pred       string `json:"pred"`
}

/*
type preds struct {
	Pred []string `json:"pred"`
}
*/
// var orders = []order{{Order_name: "Order1", Start_date: "2020-10-22"},}
var tasks = []task{
	{Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred: ""},
}

func main() {
	router := gin.Default()
	router.GET("/orders", getOrders)
	router.GET("/tasks/:id", getTasks)
	router.POST("/orders", postOrders)
	router.POST("/tasks", postTasks)
	router.PUT("/tasks", postTasks)
	fmt.Println(getOptDuration("OrderA", 10))
	router.Run("localhost:8080")
}

/*
	TODO

func selectSQL(table string) string{

}
*/

func getOptDuration(Order_name string, maxres int) float64 {
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from tasks WHERE Order_name = $1", Order_name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tasks := []task{}

	for rows.Next() {
		p := task{}
		//Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred
		err := rows.Scan(&p.Task, &p.Order_name, &p.Duration, &p.Resource, &p.Pred)
		// fmt.Println(tasks)

		if err != nil {
			fmt.Println(err)
			return -1.0
			// continue
		}
		tasks = append(tasks, p)
	}
	vartasks := tasks
	time := 0
	copy(vartasks, tasks)
	donetasks := []string{}
	// inworktasks := []string{}
	for len(vartasks) != 0 {
		num := rand.Intn(len(vartasks))
		value := vartasks[num]

		var newPreds []interface{}
		/*
			// newPreds := []string{}
			// a, _ := json.Marshal([]byte(value.Pred))
			// err := json.Unmarshal([]byte(a), &newPreds)
			fmt.Printf("Preds old: %+v \n", value.Pred)
				fmt.Printf("%T\n", value.Pred)
				fmt.Printf("%T\n", `["1","2"]`)
				fmt.Println([]byte(`["1","2"]`))
				fmt.Println([]byte(value.Pred))
				// err := json.Unmarshal([]byte(`["1","2"]`), &newPreds)
		*/
		err := json.Unmarshal([]byte(value.Pred), &newPreds)
		if err != nil {
			fmt.Println("error:", err)
		}
		// fmt.Printf("Preds: %+v \n", newPreds)
		checkPreds := func(newPreds []interface{}) bool {
			for _, i := range newPreds {
				if !inArray(i, donetasks) {
					return false
				}
			}
			return true
		}
		if !(checkPreds(newPreds)) {
			continue
		}
		// inworktasks = append(inworktasks, value.Task)
		donetasks = append(donetasks, value.Task)
		fmt.Println(donetasks)
		time += value.Duration
		vartasks = append(vartasks[:num], vartasks[num+1:]...)

		// fmt.Println(key, value)
		// if value.Pred
	}
	return float64(time)

}

// -----------ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ----------//
func inArray(val interface{}, array interface{}) (index bool) {
	values := reflect.ValueOf(array)

	if reflect.TypeOf(array).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
func RemoveIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

//END ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

// -------------------------------------------------------------------------------------
func getTasks(c *gin.Context) {
	Order_name := c.Param("id")

	// fmt.Println(Order_name)

	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from tasks WHERE Order_name = $1", Order_name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tasks := []task{}

	for rows.Next() {
		p := task{}
		//Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred
		err := rows.Scan(&p.Task, &p.Order_name, &p.Duration, &p.Resource, &p.Pred)
		// fmt.Println(tasks)

		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusConflict, err)
			continue
		}
		tasks = append(tasks, p)
	}
	// fmt.Println(tasks)
	/*
		for _, a := range tasks {
			if a.Order_name == Order_name {
				c.IndentedJSON(http.StatusOK, a)
				// return
			}
		}
	*/
	c.IndentedJSON(http.StatusOK, tasks)
}

func getOrders(c *gin.Context) {
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from orders")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	orders := []order{}

	for rows.Next() {
		p := order{}
		err := rows.Scan(&p.Order_name, &p.Start_date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, p)
	}
	// fmt.Println(orders)
	// fmt.Println(orders)
	c.IndentedJSON(http.StatusOK, orders)
	// fmt.Println(err)

}

// -----------FUNCTION postTasks---------------//
// postTasks: REST API, добавление данных по POST и PUT в таблицу tasks
func postTasks(c *gin.Context) {
	var newTask task
	//Получение данных из контекста
	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(newTask)
	//Подключение к Postgres
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	////Удаление task при обновлении
	if c.Request.Method == "PUT" {
		result, err := db.Exec("DELETE FROM tasks WHERE task = $1; ",
			newTask.Task)
		if err != nil {
			fmt.Println(result)
			c.IndentedJSON(http.StatusBadRequest, err)
			panic(err)
		}
	}
	// Добавление новой работы "task" в таблицу tasks
	result, err := db.Exec("INSERT INTO tasks (task, order_name, duration, resource, pred) values ($1, $2, $3, $4, $5)",
		newTask.Task, newTask.Order_name, newTask.Duration, newTask.Resource, newTask.Pred)
	if err != nil {
		fmt.Println(result)
		c.IndentedJSON(http.StatusBadRequest, err)
		panic(err)
	}

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

// ------------FUNCTION postOrders--------------//
// postOrders: REST API, добавление данных по POST и PUT в таблицу orders
func postOrders(c *gin.Context) {
	var newOrder order
	//Получение данных
	if err := c.BindJSON(&newOrder); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(newOrder)
	//Подключение к Postgres
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Добавление данных в таблицу
	result, err := db.Exec("insert into orders (order_name, start_date) values ($1, $2)",
		newOrder.Order_name, newOrder.Start_date)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newOrder)

}
