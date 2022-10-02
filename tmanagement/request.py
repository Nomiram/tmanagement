'''
Позволяет посылать запросы get и post на сервер
'''
## TODO: добавить полный список всех API и упорядочить
import json
import requests

URL_TASKS = 'http://localhost:8080/tasks'
URL_TASK_ORDER = 'http://localhost:8080/tasks/'+'OrderB'
URL_ORDERS = 'http://localhost:8080/orders'
URL_DURATION = 'http://localhost:8080/duration/'+'OrderB'
headers = {
    "Content-type": "application/json",
}
# Создание (обновление) Order
data = {"order_name": "OrderA", "start_date":"2020-11-23"}
print("PUT:\n",json.dumps(data))
r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data))
print(r)
print(r.text)
data = {"order_name": "OrderB", "start_date":"2020-11-23"}
print("PUT:\n",json.dumps(data))
r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data))
print(r)
print(r.text)
# Создание (обновление) Order (для удаления)
data = {"order_name": "OrderC", "start_date":"2020-11-23"}
print(json.dumps(data))
r = requests.put(URL_ORDERS, headers=headers, data=json.dumps(data))
print(r)
print(r.text)
# {Order_name: "Order1", Start_date: "2020-10-22"}
# data = {"order_name": "Order4", "start_date":"2020-11-23"}
data = {"task": "3", "order_name": "OrderA", "duration": 4, "resource": 3, "pred": '["1"]'}
print(json.dumps(data))
r = requests.put(URL_TASKS, headers=headers, data=json.dumps(data))
# r2 = requests.post(URL1, headers=headers, data=json.dumps(data))
print(r)
print(r.text)
#test1
# data = {"task": "3", "order_name": "OrderA", "duration": 4, "resource": 3, "pred": '["1"]'}
data_mas = [
{"task": "1", "order_name": "OrderA", "duration": 1, "resource": 5, "pred": '[]'},
{"task": "2", "order_name": "OrderA", "duration": 1, "resource": 5, "pred": '["1"]'},
{"task": "1", "order_name": "OrderB", "duration": 1, "resource": 5, "pred": '[]'},
{"task": "2", "order_name": "OrderB", "duration": 3, "resource": 6, "pred": '[]'},
{"task": "3", "order_name": "OrderB", "duration": 3, "resource": 4, "pred": '["1"]'},
{"task": "4", "order_name": "OrderB", "duration": 2, "resource": 3, "pred": '["1", "2"]'},
{"task": "5", "order_name": "OrderB", "duration": 10,"resource": 7, "pred": '["3"]'},
]
for data in data_mas:
    # print(json.dumps(data))
    r = requests.put(URL_TASKS, headers=headers, data=json.dumps(data))
    # r2 = requests.post(URL1, headers=headers, data=json.dumps(data))
    print(r)
    print(r.text)
r3=requests.get(URL_TASK_ORDER)
print(r3)
print(r3.text)
# Удаление Order
data = {"order_name": "OrderC"}
print(json.dumps(data))
r = requests.delete(URL_ORDERS, headers=headers, data=json.dumps(data))
print(r)
# Запрос duration
r = requests.get(URL_DURATION)
print(r)
print(r.text)
input("Press Enter to exit")
