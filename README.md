# avtoto

**avtoto** - обёртка на [API сервиса avtoto.ru](https://www.avtoto.ru/services/search/docs/technical_soap.html)

Изначально разработчики API предполагали использование SOAP-технологии, но в конечном итоге оставили [дополнение для версии на cURL](https://www.avtoto.ru/services/search/docs/technical_soap.html#curl).

Именно это дополнение используется для работы с API в этом проекте.

## Установка

```golang
go get github.com/rb-pro/avtoto
```

## С чего начать?

Для начала работы с API Вам необходимо:

- Заключите договор-поставки. Для этого обратитесь в [клиентский отдел](https://www.avtoto.ru/contacts.html) любым удобным способом. Дальнейшие действия возможны только после подписания договора.
- Активируйте сервис и добавьте IP адрес своего сайта на странице [настройка веб-сервиса](https://www.avtoto.ru/#settings:all).
- Выполните настройки на своём сайте: введите логин / пароль (как при авторизации на сайте) и номер (id) клиента (номер указан в разделе [общая информация](https://www.avtoto.ru/#settings:all)).

Вам необходимо знать:

- Номер клиента
- Логин
- Пароль

С помощью этих данных Вы можете инициилизировать пользователя:

```golang
user := avtoto.User{
    UserId:       userIdInt,
    UserLogin:    UserLoginStr,
    UserPassword: UserPasswordStr}
```

После этого Вам предоставлен функционал всего API. Методы описаны в [данной документации](https://pkg.go.dev/github.com/rb-pro/avtoto) и [документации поставщиков API](https://www.avtoto.ru/services/search/docs/technical_soap.html#curl).

Работа с данной обёрткой осуществляется с помощью работы со структурами запрос-ответ.
***Например**:* Для метода *GetBrandsByCodeRequestGetBrandsByCode существуют структуры:*

* GetBrandsByCodeRequestGetBrandsByCodeRequest - для запроса
* GetBrandsByCodeRequestGetBrandsByCodeResponse - для ответа