package main

import (
	"tmanagement/handlers"
	// "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

/*
TODO:
1. Разбить на файлы
2. Дальнейшая оптимизация кода
*/
// type handlers.Delorder = handlers.Delorder

var tasks []handlers.Task = []handlers.Task{}

// Информация для подключения к БД postgres
var CONNSTR = "user=postgres password=qwerty dbname=VS sslmode=disable"

func main() {
	router := gin.Default()
	router.GET("/duration/:handlers.Order", getBrowserOptDuration)
	router.GET("/orders", getOrders)
	router.GET("/tasks/:id", getTasks)
	router.POST("/orders", postOrders)
	router.PUT("/orders", postOrders)
	router.DELETE("/orders", delOrders)
	router.POST("/tasks", postTasks)
	router.PUT("/tasks", postTasks)
	router.DELETE("/tasks", delTasks)
	router.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"ok": "Pong"}) })
	// fmt.Println(getOptDuration("OrderB", 10, 10000))
	router.Run("localhost:8080")
}

/*
	TODO

func selectSQL(table string) string{

}
*/

/*
REST API:GET Функция возвращает кратчайшее время для работ

~/duration/<string>
*/
func getBrowserOptDuration(c *gin.Context) {
	Order_name := c.Param("handlers.Order")
	type returnstruct struct {
		Duration float64
		Path     []string
	}
	//lint:ignore SA4006 (выражение используется далее)
	path := []string{}
	i, path := getOptDuration(Order_name, 10, 100000)
	ret := returnstruct{Duration: i, Path: path}
	if i == -1 {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Status   string
			Duration float64
			Path     []string
		}{fmt.Sprint(http.StatusBadRequest), i, path})
	} else {

		c.IndentedJSON(http.StatusOK, ret)
	}
}
func getOptDuration(Order_name string, maxres int, goroutinesCount int) (float64, []string) {
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
	tasks := []handlers.Task{}
	flag := -1
	for rows.Next() {
		flag = 1
		p := handlers.Task{}
		//Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred
		err := rows.Scan(&p.Task, &p.Order_name, &p.Duration, &p.Resource, &p.Pred)

		if err != nil {
			fmt.Println(err)
			return -1.0, []string{}
			// continue
		}
		tasks = append(tasks, p)
	}
	if flag == -1 {
		return -1.0, []string{}
	}
	tasksEn := map[string]handlers.TaskEn{}
	for _, tas := range tasks {
		var newPreds []string
		err := json.Unmarshal([]byte(tas.Pred), &newPreds)
		if err != nil {
			fmt.Println("Cannot Unmarshal")
			return -1, []string{}
		}

		tasksEn[tas.Task] = handlers.TaskEn{Task: tas.Task, Order_name: tas.Order_name, Duration: tas.Duration, Resource: tas.Resource, Pred: newPreds}
	}

	type ret struct {
		Duration float64
		Path     []string
	}

	doCh := make(chan ret)
	GPSS := func() {
		// Начало имитации работы для вычисления длительности проекта
		vartasks := tasksEn // Массив данных для работ

		waitingtasks := []string{}      // Список оставшихся работ
		var posibletasks = []string{}   // Список допустимых для выполнения работ(предыдущие работы завершены)
		time := 0                       // Время выполнения
		donetasks := []string{}         // Список работ, которые были завершены
		inworktasks := map[string]int{} // Список работ, которые сейчас выполняются

		// Формирование списка работ
		for k := range vartasks {
			waitingtasks = append(waitingtasks, k)
		}
		// Пока остались незавершенные работы
		for len(waitingtasks) > 0 || len(inworktasks) > 0 {
			for {
				// Подсчет суммы занятых ресурсов
				sum_ := 0
				for key := range inworktasks {
					sum_ += vartasks[key].Resource
				}
				//Формирование массива возможных работ
				posibletasks = []string{}
				for _, tas := range waitingtasks {

					// Проверка: готовы ли обязательные предыдущие работы
					checkPreds := func(value handlers.TaskEn) bool {
						for _, i := range value.Pred {
							if !inArray(i, donetasks) {
								return false
							}
						}
						return true
					}
					if !(checkPreds(vartasks[tas])) {
						continue
					}

					//Проверка доступности по ресурсам
					if sum_+vartasks[tas].Resource > maxres {
						continue
					}
					posibletasks = append(posibletasks, tas)
				}
				// Если есть доступные работы
				if len(posibletasks) > 0 {
					// Выбираем случайную работу из списка доступных
					num := rand.Intn(len(posibletasks))
					value := vartasks[posibletasks[num]]
					// Добавляем работу в массив выполняющихся работ
					inworktasks[value.Task] = value.Duration
					// Удаляем из массива возможных
					for ind, val := range waitingtasks {
						if val == posibletasks[num] {
							waitingtasks = append(waitingtasks[:ind], waitingtasks[ind+1:]...)
							break
						}
					}
					// Удаление значения из posibletasks
					//lint:ignore SA4006 (выражение используется далее)
					posibletasks = append(posibletasks[:num], posibletasks[num+1:]...)
				} else {
					break
				}
			}
			//Если не было добавлено ничего и ничего не осталось, то
			if len(inworktasks) == 0 {
				if len(waitingtasks) > 0 {
					doCh <- ret{-1.0, donetasks}
				} else {
					// return float64(time)
					doCh <- ret{float64(time), donetasks}

				}
				return
			}
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
		}
		// return float64(time)
		doCh <- ret{float64(time), donetasks}
		// return

	}
	start := time.Now() //Запись времени

	for i := 0; i < goroutinesCount; i++ { //Запуск goroutinesCount горутин
		go GPSS()

	}
	mintime := -1.0
	mas := []ret{}
	for i := 0; i < goroutinesCount; i++ {
		mas = append(mas, <-doCh)
	}
	unique := map[float64]bool{}

	for _, v := range mas {
		unique[v.Duration] = true
	}

	fmt.Print("unique: ")

	for key := range unique {
		fmt.Print(key, " ")
	}
	fmt.Println()
	mintime = mas[0].Duration
	minpath := []string{}
	for i := 0; i < goroutinesCount; i++ {
		if mas[i].Duration <= mintime {
			mintime = mas[i].Duration
			minpath = mas[i].Path
		}
	}
	duration2 := time.Since(start)
	fmt.Println("Время с параллелизмом: ", duration2, "путь", minpath, "количество горутин: ", len(mas))
	return float64(mintime), minpath
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

// Возвращает минимум двух чисел int
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

/*
// Возвращает минимум двух чисел float64
func minfl(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}
*/
//END ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// -------------------------------------------------------------------------------------

// REST API:GET Возвращает информацию из таблицы tasks
func getTasks(c *gin.Context) {
	Order_name := c.Param("id")

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
	tasks := []handlers.Task{}

	for rows.Next() {
		p := handlers.Task{}
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
	orders := []handlers.Order{}

	for rows.Next() {
		p := handlers.Order{}
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
func postTasks(c *gin.Context) {
	var newTask handlers.Task
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
	////Удаление handlers.Task при обновлении
	if c.Request.Method == "PUT" {
		result, err := db.Exec("DELETE FROM tasks WHERE handlers.Task = $1 AND order_name = $2; ",
			newTask.Task, newTask.Order_name)
		if err != nil {
			fmt.Println(result)
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
	}
	// Добавление новой работы "handlers.Task" в таблицу tasks
	result, err := db.Exec("INSERT INTO tasks (handlers.Task, order_name, duration, resource, pred) values ($1, $2, $3, $4, $5)",
		newTask.Task, newTask.Order_name, newTask.Duration, newTask.Resource, newTask.Pred)
	if err != nil {
		fmt.Println(result)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

// ------------FUNCTION postOrders--------------//

// REST API:POST,PUT добавление данных по POST и PUT в таблицу orders
//
// json: {"order_name":string, start_date":string}
func postOrders(c *gin.Context) {
	var newOrder handlers.Order
	//Получение данных
	if err := c.BindJSON(&newOrder); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	////Удаление handlers.Order при обновлении
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

	c.IndentedJSON(http.StatusCreated, newOrder)

}

// REST API:DELETE удаление данных из таблицы orders
//
// json: {"order_name":string}
func delOrders(c *gin.Context) {
	var Order handlers.Delorder
	//Получение данных
	if err := c.BindJSON(&Order); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", CONNSTR)
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

	c.IndentedJSON(http.StatusCreated, Order)

}

// REST API:DELETE удаление данных из таблицы orders
//
// json: {"order_name":string}
func delTasks(c *gin.Context) {
	var delTask handlers.TaskDel
	//Получение данных
	if err := c.BindJSON(&delTask); err != nil {
		fmt.Println(err)
		return
	}
	//Подключение к Postgres
	db, err := sql.Open("postgres", CONNSTR)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Удаление данных из таблицы
	result, err := db.Exec("DELETE FROM tasks WHERE order_name = $1 AND handlers.Task = $2",
		delTask.Order_name, delTask.Task)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, delTask)

}
