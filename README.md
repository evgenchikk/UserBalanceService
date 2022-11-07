### Launch notes

> Локальный запуск

В каталоге с файлом docker-compose.yml выполнить команду
- docker-compose up

> Запуск на сервере

В файле docker-compose.yml изменить значения в секции environment <br/>&emsp;HOST: \<server ip\>

В каталоге с файлом docker-compose.yml выполнить команду
- docker-compose up

По умолчанию сервис доступен на порте 3000. Порт можно изменить в файле docker-compose.yml в секции environment, при этом следует еще поправить проброс порта в секции ports. База данных доступна на порте 10330.

---
### API

> *Swagger*

- UI -- /docs/index.html
- json -- /docs/doc.json


---
> *Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.*

- POST /add <br/>
{ <br/>
    &emsp;"user_id": user_id, <br/>
    &emsp;"amount": amount <br/>
}

curl -d '{"user_id": 1, "amount": 100}' -X POST localhost:3000/add


---
> *Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость.*

- POST /reserve <br/>
{ <br/>
    &emsp;"user_id": user_id, <br/>
    &emsp;"service_id": service_id, <br/>
    &emsp;"order_id": order_id, <br/>
    &emsp;"price": price <br/>
}

curl -d '{"user_id": 1, "service_id": 1, "order_id": 1, "price": 1000}' -X POST localhost:3000/reserve


---
> *Метод разрезервирования средств пользователя. Принимает id заказа.*

- POST /dereserve <br/>
{ <br/>
    &emsp;"order_id": order_id <br/>
}

curl -d '{"order_id": 1}' -X POST localhost:3000/dereserve


---
> *Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму.*

- POST /approve <br/>
{ <br/>
    &emsp;"user_id": user_id, <br/>
    &emsp;"service_id": service_id, <br/>
    &emsp;"order_id": order_id, <br/>
    &emsp;"amount": amount <br/>
}

curl -d '{"user_id": 1, "service_id": 1, "order_id": 1, "amount": 1000}' -X POST localhost:3000/approve


---
> *Метод перевода денег с баланса одного пользователя на баланс другого. Принимает id пользователя-отправителя, id пользователя-получателя, сумму.*

- POST /transfer <br/>
{ <br/>
    &emsp;"from_user_id": user_id, <br/>
    &emsp;"to_user_id": user_id, <br/>
    &emsp;"amount": amount <br/>
}

curl -d '{"from_user_id": 1, "to_user_id": 2, "amount": 100}' -X POST localhost:3000/transfer


---
> *Метод получения баланса пользователя. Принимает id пользователя.*

- POST /balance <br/>
{ <br/>
    &emsp;"user_id": user_id <br/>
}

curl -d '{"user_id": 1}' -X POST localhost:3000/balance


---
> *Метод получения отчета. Принимает год-месяц.*

- POST /report <br/>
{ <br/>
    &emsp;"period": "yyyy-mm" <br/>
}

curl -d '{"period": "2022-10"}' -X POST localhost:3000/report


---
> *Скачивание отчета. Принимает имя файла отчета, полуенного в ответе на запрос создания отчета.*

- GET /report/\<report filename\>

curl -X GET localhost:3000/report/\<report filename\>


---
### Модель базы данных
<img alt="Database model" src="/database/db_model.png">

- <b>users</b><br/>
Таблица пользователей. Хранит id пользователя, его баланс и зарезервированные деньги (reserved_balance).
- <b>orders</b><br/>
Таблица заказов услуг (service). Хранит id заказа, id пользователя, id услуги, время заказа и цену.
- <b>payment_history</b><br/>
Таблица истории платежей. Хранит id заказа, время платежа и сумму платежа. Используется для получения отчетов для бухгалтерии. Предполагается, что каждый пользователь может оплатить только свои заказы, поэтому в этой таблице не содержится атрибут user_id.


---
### Комментарии

Маршрут /balance работает через http метод POST и принимает id пользователя через тело запроса. Возможно, стоит включить id пользователя в сам маршрут (т.е. /balance/\<user_id\>). Но я решил, руководствуясь описанием задания ("Сервис должен предоставлять HTTP API с форматом JSON как при отправке запроса, так и при получении результата"), что полностью все взаимодействие с сервисом должно быть через JSON.

Проверка баланса на неотрицательные значения выполняется на уровне базы данных - установлен check на колонках с балансами.

Отменить заказ (/dereserve) можно, только если не было признано выручки по данному заказу. Т. е. если по данному заказу отсутствуют записи в таблице payment_history.
