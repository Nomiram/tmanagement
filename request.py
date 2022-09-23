'''
Позволяет посылать запросы get и post на сервер
'''
# TODO: добавить полный список всех API и упорядочить
import json
import requests

URL1 = 'http://localhost:8080/tasks'
URL2 = 'http://localhost:8080/tasks/OrderA'
URL3 = 'http://localhost:8080/orders'
headers = {
    "Content-type": "application/json",
}
# {Order_name: "Order1", Start_date: "2020-10-22"}
# data = {"order_name": "Order4", "start_date":"2020-11-23"}
data = {"task": "3", "order_name": "OrderA", "duration": 4, "resource": 3, "pred": "1"}
print(json.dumps(data))
r2 = requests.put(URL1, headers=headers, data=json.dumps(data))
# r2 = requests.post(URL1, headers=headers, data=json.dumps(data))
print(r2)
print(r2.text)
r3=requests.get(URL2)
print(r3)
print(r3.text)
