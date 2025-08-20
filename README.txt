                ТРЕБОВАНИЯ
                ``````````

GET запрос на localhost:8080/calculations
возвращает историю вычислений

POST запрос на localhost:8080/calculations
передает json с expression

PATCH запрос на localhost:8080/calculations/id 
телом запросом формата json передаем новый expression

DELETE запрос на localhost:8080/calculations/id 
передаем id (без тела запроса)