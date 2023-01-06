// avtoto - обёртка на API сервиса [avtoto.ru]
//
// Изначально разработчики API предполагали использование SOAP-технологии, но в конечном итоге оставили [дополнение для версии на cURL].
// Именно это дополнение используется для работы с API в этом проекте.
//
// # Установка
//
//	go get github.com/rb-pro/avtoto
//
// # С чего начать?
//
// Для начала работы с API Вам необходимо:
//   - Заключите договор-поставки. Для этого обратитесь в [клиентский отдел] любым удобным способом. Дальнейшие действия возможны только после подписания договора.
//   - Активируйте сервис и добавьте IP адрес своего сайта на странице [настройка веб-сервиса].
//   - Выполните настройки на своём сайте: введите логин / пароль (как при авторизации на сайте) и номер (id) клиента (номер указан в разделе [общая информация]).
//
// Вам необходимо знать:
//   - Номер клиента
//   - Логин
//   - Пароль
//
// С помощью этих данных Вы можете инициилизировать пользователя:
//
//	user := avtoto.User{
//		UserId:       userIdInt,
//		UserLogin:    UserLoginStr,
//		UserPassword: UserPasswordStr}
//
// После этого Вам предоставлен функционал всего API. Методы описаны в [данной документации] и [документации поставщиков API].
//
// [avtoto.ru]: https://www.avtoto.ru/services/search/docs/technical_soap.html
// [дополнение для версии на cURL]: https://www.avtoto.ru/services/search/docs/technical_soap.html#curl
// [клиентский отдел]: https://www.avtoto.ru/contacts.html
// [настройка веб-сервиса]: https://www.avtoto.ru/#settings:all
// [общая информация]: https://www.avtoto.ru/#settings:all
// [данной документации]: https://pkg.go.dev/github.com/rb-pro/avtoto
// [документации поставщиков API]: https://www.avtoto.ru/services/search/docs/technical_soap.html#curl
package avtotoGo

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

const URL string = "https://www.avtoto.ru/?soap_server=json_mode"

// Исходная структура для авторизации пользователя
type User struct {
	UserId       int    `json:"user_id"`       // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
}

// Запрос с параметром action и данными json в формате []byte
func httpPost(bytesRepresentation []byte, action string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("action", action)
	_ = writer.WriteField("data", string(bytesRepresentation))
	responseError := writer.Close()
	if responseError != nil {
		return nil, responseError
	}

	client := &http.Client{}
	req, responseError := http.NewRequest(http.MethodPost, URL, payload)
	if responseError != nil {
		return nil, responseError
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, responseError := client.Do(req)
	if responseError != nil {
		return nil, responseError
	}
	defer res.Body.Close()

	// Считываем ответ
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, responseError
		}
		return bodyBytes, responseError
	} else {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
}
