'''
Позволяет посылать запросы get и post на сервер
'''
## TODO: добавить полный список всех API и упорядочить
import json
import time

import requests

URL_TASKS = 'http://localhost:8080/tasks'
URL_TASK_ORDER = 'http://localhost:8080/tasks/'+'OrderB'
URL_TASK_ORDERA = 'http://localhost:8080/tasks/'+'OrderA'
URL_TASK_ORDERB = 'http://localhost:8080/tasks/'+'OrderB'
# URL_TASK_ORDER = 'http://localhost:8080/tasks/'+'OrderB'
URL_ORDERS = 'http://localhost:8080/orders'
URL_DURATION = 'http://localhost:8080/duration/'+'OrderB'
headers = {
    "Content-type": "application/json",
}
def testing():
    print("Тестирование:")
    print("Создание (обновление) Order: ",end="")
    # Создание (обновление) Order
    data = {"order_name": "OrderA", "start_date":"2020-11-23"}
    # print("PUT:\n",json.dumps(data))
    r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data),timeout=10)
    if 200<=r.status_code<300: print("test passed")
    else:
        print("test failed\n", r.text)
        return 1
    data = {"order_name": "OrderB", "start_date":"2020-11-23"}
    # print("PUT:\n",json.dumps(data))
    r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data),timeout=10)
    if 200<=r.status_code<300: print("test#2 passed")
    else: print("test#2 failed\n", r.text)
    print("Создание (обновление) OrderC (для удаления): ",end="")
    # Создание (обновление) Order (для удаления)
    data = {"order_name": "OrderC", "start_date":"2020-11-23"}
    # print(json.dumps(data))
    r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data),timeout=10)
    if 200<=r.status_code<300: print("test#3 passed")
    else: print("test#3 failed\n", r.text)
    # Удаление Order
    print("Удаление OrderC: ",end="")
    data = {"order_name": "OrderC"}
    # print(json.dumps(data))
    r = requests.delete(URL_ORDERS, headers=headers, data=json.dumps(data),timeout=10)
    if 200<=r.status_code<300: print("delete test passed")
    else: print("delete test failed\n", r.text)

    print("Добавление (обновление) списка tasks для расчетов: ",end="")
    # {Order_name: "Order1", Start_date: "2020-10-22"}
    # data = {"order_name": "Order4", "start_date":"2020-11-23"}
    data = {"task": "3", "order_name": "OrderA", "duration": 4, "resource": 3, "pred": '["1"]'}
    # print(json.dumps(data))
    r = requests.put(URL_TASKS, headers=headers, data=json.dumps(data),timeout=10)
    # r2 = requests.post(URL1, headers=headers, data=json.dumps(data))
    if 200<=r.status_code<300: print("put test passed")
    else: print("put test failed\n", r.text)
    #test1
    # data = {"task": "3", "order_name": "OrderA", "duration": 4, "resource": 3, "pred": '["1"]'}
    data_mas = [
    {"task": "10", "order_name": "OrderA", "duration": 1, "resource": 5, "pred": '[]'},
    {"task": "20", "order_name": "OrderA", "duration": 1, "resource": 5, "pred": '["1"]'},
    {"task": "1", "order_name": "OrderB", "duration": 1, "resource": 5, "pred": '[]'},
    {"task": "2", "order_name": "OrderB", "duration": 3, "resource": 6, "pred": '[]'},
    {"task": "3", "order_name": "OrderB", "duration": 3, "resource": 4, "pred": '["1"]'},
    {"task": "4", "order_name": "OrderB", "duration": 2, "resource": 3, "pred": '["1", "2"]'},
    {"task": "5", "order_name": "OrderB", "duration": 10,"resource": 7, "pred": '["3"]'},
    ]
    for data in data_mas:
        # print(json.dumps(data))
        r = requests.put(URL_TASKS, headers=headers, data=json.dumps(data),timeout=10)
        # r2 = requests.post(URL1, headers=headers, data=json.dumps(data))
        if 200<=r.status_code<300: pass
        else: print("put test failed\n", r.text)
    print("Проверка добавления: ",end="")
    r3=requests.get(URL_TASK_ORDERA,timeout=10)
    testarr = json.loads(r3.text)
    r3=requests.get(URL_TASK_ORDERB,timeout=10)
    testarr += json.loads(r3.text)
    if 200<=r3.status_code<300:
        for i in data_mas:
            if i in testarr:
                pass
                # print("pass!")
            else:
                print("err: ", i)
                print("get test failed\n", r3.text)
                # input("Press Enter to exit")
                return
        print("get test passed")
        # print(r3.text)
    else: print("get test failed\n", r3.text)

    print("Запрос duration")
    start_time = time.time()
    # Запрос duration
    r = requests.get(URL_DURATION,timeout=10)
    if 200<=r.status_code<300: print("duration test passed")
    else: print("duration test failed\n", r.text)
    print(r.text)
    print(f"{(time.time() - start_time)} seconds")
    # Запрос duration повторный
    print("Повторный")
    start_time = time.time()
    r = requests.get(URL_DURATION,timeout=10)
    if 200<=r.status_code<300: print("duration test passed")
    else: print("duration test failed\n", r.text)
    print(r.text)
    print(f"{(time.time() - start_time)} seconds")
    print("Добавление нового элемента: ",end="")
    data = {"task": "6", "order_name": "OrderB", "duration": 4, "resource": 10, "pred": '["1"]'}
    r = requests.put(URL_TASKS, headers=headers, data=json.dumps(data),timeout=10)
    if 200<=r.status_code<300: 
        print("put test passed")
        # print(r.text)
    else: print("put test failed\n", r.text)
    print("Повторный с добавлением: ",end="")
    start_time = time.time()
    # Запрос duration повторный2
    r = requests.get(URL_DURATION,timeout=10)
    if 200<=r.status_code<300: print("duration test passed")
    else: print("duration test failed\n", r.text)
    print(r.text)
    print(f"{(time.time() - start_time)} seconds")
    print("Повторный#2 (без изменений): ",end="")
    start_time = time.time()
    # Запрос duration повторный2
    r = requests.get(URL_DURATION,timeout=10)
    if 200<=r.status_code<300: print("duration test passed")
    else: print("duration test failed\n", r.text)
    print(r.text)
    print(f"{(time.time() - start_time)} seconds")
    # Удаление Task
    data = {"order_name": "OrderB", "task": "6"}
    print(json.dumps(data))
    r = requests.delete(URL_TASKS, headers=headers, data=json.dumps(data),timeout=10)
    return 0

try:
    if testing() == 0:
        print("all tests passed!")
    input("Press Enter to exit")
except Exception as e:
    print("Exception:", e.with_traceback(e.__traceback__))
    input("Press Enter to exit")
