package main

import (
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
2. Добавить удаление REST API DELETE
3. Дальнейшая оптимизация кода
4. Возвращение упорядоченного списка задач в дополнение к длительности
*/

// Информация для подключения к БД postgres
var CONNSTR = "user=postgres password=qwerty dbname=VS sslmode=disable"

type order struct {
	Order_name string `json:"order_name"`
	Start_date string `json:"start_date"`
}
type delorder struct {
	Order_name string `json:"order_name"`
}
type task struct {
	Task       string `json:"task"`
	Order_name string `json:"order_name"`
	Duration   int    `json:"duration"`
	Resource   int    `json:"resource"`
	Pred       string `json:"pred"`
}

type taskEn struct {
	Task       string   `json:"task"`
	Order_name string   `json:"order_name"`
	Duration   int      `json:"duration"`
	Resource   int      `json:"resource"`
	Pred       []string `json:"pred"`
}
type taskDel struct {
	Task       string `json:"task"`
	Order_name string `json:"order_name"`
}

/*
type preds struct {
	Pred []string `json:"pred"`
}
*/
// var orders = []order{{Order_name: "Order1", Start_date: "2020-10-22"},}

var tasks = []task{
	{Task: "1", Order_name: "Order1", Duration: 2, Resource: 3, Pred: "[]"},
}

func main() {
	router := gin.Default()
	router.GET("/duration/:order", getBrowserOptDuration)
	router.GET("/orders", getOrders)
	router.GET("/tasks/:id", getTasks)
	router.POST("/orders", postOrders)
	router.PUT("/orders", postOrders)
	router.DELETE("/orders", delOrders)
	router.POST("/tasks", postTasks)
	router.PUT("/tasks", postTasks)
	router.DELETE("/tasks", delTasks)
	fmt.Println(getOptDuration("OrderB", 10, 10000))
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
	Order_name := c.Param("order")
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
			return -1.0, []string{}
			// continue
		}
		tasks = append(tasks, p)
	}
	if flag == -1 {
		return -1.0, []string{}
	}
	/*
		type wtask struct {
			Task    string
			RemTime int
		}
	*/
	tasksEn := map[string]taskEn{}
	for _, tas := range tasks {
		var newPreds []string
		err := json.Unmarshal([]byte(tas.Pred), &newPreds)
		if err != nil {
			fmt.Println("Cannot Unmarshal")
			return -1, []string{}
		}

		tasksEn[tas.Task] = taskEn{tas.Task, tas.Order_name, tas.Duration, tas.Resource, newPreds}
	}
	type ret struct {
		Duration float64
		Path     []string
	}

	doCh := make(chan ret)
	GPSS := func() {
		// Начало имитации работы для вычисления длительности проекта
		vartasks := tasksEn // Массив данных для работ

		waitingtasks := []string{}    // Список оставшихся работ
		var posibletasks = []string{} // Список допустимых для выполнения работ(предыдущие работы завершены)
		// copy(vartasks, tasksEn)
		time := 0                       // Время выполнения
		donetasks := []string{}         // Список работ, которые были завершены
		inworktasks := map[string]int{} // Список работ, которые сейчас выполняются
		// var value task
		// var num int

		// Формирование списка работ
		for k := range vartasks {
			waitingtasks = append(waitingtasks, k)
		}
		// fmt.Println(len(vartasks))
		// fmt.Println(waitingtasks)
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
					checkPreds := func(value taskEn) bool {
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
					// fmt.Println("undone")
					doCh <- ret{-1.0, donetasks}
				} else {
					// return float64(time)
					// fmt.Println("done:", donetasks)
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
					// fmt.Println()
					donetasks = append(donetasks, ind)
					// fmt.Println("testd:", donetasks, ind)
					delete(inworktasks, ind)
				}
			}
			time += mintime
			//
			// return -1
			// RemoveIndex(&vartasks, num)
			//Удаление индекса num
			// vartasks = append(vartasks[:num], vartasks[num+1:]...)

			// fmt.Println(key, value)
			// if value.Pred
		}
		// return float64(time)
		// fmt.Println("done:", donetasks)
		doCh <- ret{float64(time), donetasks}
		// return

	}
	start := time.Now() //Запись времени

	for i := 0; i < goroutinesCount; i++ { //Запуск goroutinesCount горутин
		go GPSS()

	}
	mintime := -1.0
	mas := []ret{}
	/*
		result := make(chan float64)
		go func ()  {

		}
	*/
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

// REST API:POST,PUT добавление данных по POST и PUT в таблицу tasks
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
		result, err := db.Exec("DELETE FROM tasks WHERE task = $1 AND order_name = $2; ",
			newTask.Task, newTask.Order_name)
		if err != nil {
			fmt.Println(result)
			c.IndentedJSON(http.StatusBadRequest, err)
			// panic(err)
			return
		}
	}
	// Добавление новой работы "task" в таблицу tasks
	result, err := db.Exec("INSERT INTO tasks (task, order_name, duration, resource, pred) values ($1, $2, $3, $4, $5)",
		newTask.Task, newTask.Order_name, newTask.Duration, newTask.Resource, newTask.Pred)
	if err != nil {
		fmt.Println(result)
		c.IndentedJSON(http.StatusBadRequest, err)
		// panic(err)
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
	////Удаление order при обновлении
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
	var Order delorder
	//Получение данных
	if err := c.BindJSON(&Order); err != nil {
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
	var delTask taskDel
	//Получение данных
	if err := c.BindJSON(&delTask); err != nil {
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
	//Удаление данных из таблицы
	result, err := db.Exec("DELETE FROM tasks WHERE order_name = $1 AND task = $2",
		delTask.Order_name, delTask.Task)
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, delTask)

}
