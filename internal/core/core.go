package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
	"tmanagement/internal/headers"
)

func GetOptDuration(Order_name string, maxres int, goroutinesCount int) (float64, []string) {
	// Получение всех работ для задачи
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
	flag := -1
	for rows.Next() {
		flag = 1
		p := headers.Task{}
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
	tasksEn := map[string]headers.TaskEn{}
	for _, tas := range tasks {
		var newPreds []string
		err := json.Unmarshal([]byte(tas.Pred), &newPreds)
		if err != nil {
			fmt.Println("Cannot Unmarshal")
			return -1, []string{}
		}

		tasksEn[tas.Task] = headers.TaskEn{Task: tas.Task, Order_name: tas.Order_name, Duration: tas.Duration, Resource: tas.Resource, Pred: newPreds}
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
					checkPreds := func(value headers.TaskEn) bool {
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
