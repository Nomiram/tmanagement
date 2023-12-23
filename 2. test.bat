REM Set-PSDebug -Trace 1
prompt ^> 
chcp 1251
@echo  "Ping Test:"
curl.exe -X GET http://localhost:8080/ping -H "Content-type: application/json"
@echo.
@echo  "Создание (обновление) Order:"
curl.exe -X PUT http://localhost:8080/orders -H "Content-type: application/json" -d "{""order_name"": ""OrderA"", ""start_date"": ""2020-11-23""}"
@echo.
@echo  "Создание (обновление) списка tasks для расчетов:"
curl.exe -X PUT http://localhost:8080/tasks -H "Content-type: application/json" -d "{""task"": ""3"", ""order_name"": ""OrderA"", ""duration"": 4, ""resource"": 3, ""pred"": ""[\""1\""]""}"
@echo.
@echo  "Проверка добавления:"
curl.exe -X GET http://localhost:8080/tasks/OrderA
@echo.
@echo  "Запрос duration:"
curl.exe -X GET http://localhost:8080/duration/OrderB
@echo.
@echo  "Добавление нового элемента:"
curl.exe -X PUT http://localhost:8080/tasks -H "Content-type: application/json" -d "{""task"": ""6"", ""order_name"": ""OrderB"", ""duration"": 4, ""resource"": 10, ""pred"": ""[\""1\""]""}"
@echo.
@echo  "Запрос duration:"
curl.exe -X GET http://localhost:8080/duration/OrderB
@echo  "Удаление Task:"
curl.exe -X DELETE http://localhost:8080/tasks -H "Content-type: application/json" -d "{""order_name"": ""OrderB"", ""task"": ""6""}"

pause