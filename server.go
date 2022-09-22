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
	router.GET("/duration/:order", getBrowserOptDuration)
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
func getBrowserOptDuration(c *gin.Context) {
	Order_name := c.Param("order")
	type returnstruct struct{ Duration float64 }
	i := getOptDuration(Order_name, 10)
	ret := returnstruct{Duration: i}
	if i == -1 {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Status   string
			Duration float64
		}{fmt.Sprint(http.StatusBadRequest), i})
	} else {

		c.IndentedJSON(http.StatusOK, ret)
	}
}
func getOptDuration(Order_name string, maxres int) float64 {
	// Получение всех работ для задачи
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
	flag := -1
	for rows.Next() {
		flag = 1
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
	if flag == -1 {
		return -1.0
	}
	/*
		type wtask struct {
			Task    string
			RemTime int
		}
	*/
	// Начало имитации работы для вычисления длительности проекта
	vartasks := tasks             // Массив данных для работ
	waitingtasks := []string{}    // Список оставшихся работ
	var posibletasks = []string{} // Список допустимых для выполнения работ(предыдущие работы завершены)
	copy(vartasks, tasks)
	time := 0                       // Время выполнения
	donetasks := []string{}         // Список работ, которые были завершены
	inworktasks := map[string]int{} // Список работ, которые сейчас выполняются
	// var value task
	// var num int

	// Формирование списка работ
	for _, tas := range vartasks {
		waitingtasks = append(waitingtasks, tas.Task)
	}
	// fmt.Println(vartasks)
	// fmt.Println(waitingtasks)
	// Пока остались незавершенные работы
	for len(waitingtasks) > 0 || len(inworktasks) > 0 {
		for {
			//Формирование массива возможных работ
			posibletasks = []string{}
			for _, value := range vartasks {
				if inArray(value.Task, waitingtasks) {

					// Проверка: готовы ли обязательные предыдущие работы
					var newPreds []interface{}
					err := json.Unmarshal([]byte(value.Pred), &newPreds)
					if err != nil {
						fmt.Println("error:", err)
					}
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

					// fmt.Println("ps1 ", posibletasks)
					//Проверка доступности по ресурсам
					sumResource := func(tasks []task, active map[string]int) int {
						sum_ := 0
						for _, tas := range tasks {
							for key := range active {
								if tas.Task == key {
									sum_ += tas.Resource
								}
							}
						}
						return sum_
					}(vartasks, inworktasks)
					if sumResource+value.Resource > maxres {
						continue
					}
					//добавляем работу в массив возможных работ
					posibletasks = append(posibletasks, value.Task)
					// fmt.Println("ps2 ", posibletasks)
				}
			}
			// Если есть доступные работы
			// return -1
			if len(posibletasks) > 0 {
				// Выбираем случайную работу из списка доступных
				num := rand.Intn(len(posibletasks))
				value := task{}
				for _, val := range vartasks {
					if val.Task == posibletasks[num] {
						value = val
					}
				}
				//добавляем работу в массив выполняющихся работ
				// inworktasks = append(inworktasks, wtask{Task: value.Task, RemTime: value.Duration})
				inworktasks[value.Task] = value.Duration
				// Удаляем из массива возможных
				for ind, val := range waitingtasks {
					if val == posibletasks[num] {
						waitingtasks = append(waitingtasks[:ind], waitingtasks[ind+1:]...)
						break
					}
				}
				// Удаление значения из posibletasks
				posibletasks = append(posibletasks[:num], posibletasks[num+1:]...)
				// fmt.Println("posi ", posibletasks)

			} else {
				break
			}
			// fmt.Println("inv33 ", inworktasks)
		}
		// return -1
		//Если не было добавлено ничего и ничего не осталось, то
		if len(inworktasks) == 0 {
			return float64(time)
		}
		// fmt.Println("inv ", inworktasks)
		// fmt.Println("ps ", posibletasks)
		//Переход к следующему времени
		mintime := 0
		for _, Remtime := range inworktasks {
			if mintime > 0 {
				mintime = min(mintime, Remtime)
			} else {
				mintime = Remtime
			}
		}
		for ind := range inworktasks {
			inworktasks[ind] -= mintime
			if inworktasks[ind] <= 0 {
				donetasks = append(donetasks, ind)
				delete(inworktasks, ind)
			}
		}
		time += mintime
		//
		// fmt.Println(donetasks)
		// return -1
		// RemoveIndex(&vartasks, num)
		//Удаление индекса num
		// vartasks = append(vartasks[:num], vartasks[num+1:]...)

		// fmt.Println(key, value)
		// if value.Pred
	}
	return float64(time)

}

// -----------ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ----------//

// Проверяет есть ли val в массиве array
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

// Возвращает минимум двух чисел
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Not working
func RemoveIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

//END ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// -------------------------------------------------------------------------------------

// Возвращает ответ на REST API запрос
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
