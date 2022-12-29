package avtotoGo

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const URL string = "https://www.avtoto.ru/?soap_server=json_mode"

// Исходная структура для авторизации пользователя
type User struct {
	UserId       int    `json:"user_id"`       // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
}

// Запрос с параметром action и данными json в формате []byte
func HttpPost(bytesRepresentation []byte, action string) ([]byte, error) {
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
	body, responseError := ioutil.ReadAll(res.Body)
	if responseError != nil {
		return nil, responseError
	}
	return body, responseError
}
