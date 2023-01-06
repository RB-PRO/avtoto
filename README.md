# avtoto

[![Go Reference](https://pkg.go.dev/badge/github.com/rb-pro/avtoto.svg)](https://pkg.go.dev/github.com/rb-pro/avtoto) [![avtoto API](https://img.shields.io/badge/avtoto-API-blue.svg)](https://www.avtoto.ru/services/search/docs/technical_soap.html)

<img align="right" alt="DiscordGo logo" src="docs/img/avtotoGO_rectangle.png" width="200">

**avtoto** - обёртка на [API сервиса avtoto.ru](https://www.avtoto.ru/services/search/docs/technical_soap.html)

Изначально разработчики API предполагали использование SOAP-технологии, но в конечном итоге оставили [дополнение для версии на cURL](https://www.avtoto.ru/services/search/docs/technical_soap.html#curl). Именно это дополнение используется для работы с API в этом проекте.

## Установка

```sh
go get github.com/rb-pro/avtoto
```

## С чего начать?

Для начала работы с API Вам необходимо:

- Заключить договор-поставки. Для этого обратитесь в [клиентский отдел](https://www.avtoto.ru/contacts.html) любым удобным способом. Дальнейшие действия возможны только после подписания договора.
- Активировать сервис и добавьте IP адрес своего сайта на странице [настройка веб-сервиса](https://www.avtoto.ru/#settings:all).
- Ввести логин / пароль (как при авторизации на сайте) и номер (id) клиента (номер указан в разделе [общая информация](https://www.avtoto.ru/#settings:all)).

Вам необходимо знать:

- Номер клиента
- Логин
- Пароль

С помощью этих данных Вы можете инициализировать пользователя:

```golang
user := avtoto.User{
    UserId:       userIdInt,
    UserLogin:    UserLoginStr,
    UserPassword: UserPasswordStr}
```

После этого Вам предоставлен функционал всего API. Методы описаны в [данной документации](https://pkg.go.dev/github.com/rb-pro/avtoto) и [документации поставщиков API](https://www.avtoto.ru/services/search/docs/technical_soap.html#curl).

Работа с данной обёрткой осуществляется с помощью работы со структурами запрос-ответ.
***Например**:* Для метода *GetBrandsByCodeRequestGetBrandsByCode* существуют *структуры:*

* GetBrandsByCodeRequestGetBrandsByCodeRequest - для запроса
* GetBrandsByCodeRequestGetBrandsByCodeResponse - для ответа
