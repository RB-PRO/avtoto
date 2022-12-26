package main

import (
	"encoding/json"
	"errors"
)

// Запрос
type SearchStartRequestStruct struct {
	UserId       int    `json:"user_id"`         // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`      // Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"`   // Пароль пользователя (тип: строка)
	SearchCode   string `json:"search_code"`     // Поисковый запрос, минимум 3 символа (тип: строка)
	SearchCross  string `json:"search_cross"`    // Искать в аналогах или нет (тип: строка, 'on' или 'off')
	Brand        string `json:"brand,omitempty"` /// Искать код с учетом бренда, минимум 2 символа (опционально)(тип: строка)
	// [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
	// [**] Список брендов можно получить с помощью метода GetBrandsByCode
}

// Ответ
type SearchStartResponseStruct struct {
	ProcessSearchID string `json:"ProcessSearchId"` // идентификатор процесса поиска (тип: строка). Необходим для отслеживания результатов процесса поиска.
	Info            struct {
		SearchID string   `json:"SearchId"`
		Errors   []string `json:"Errors"`
		Logs     string   `json:"Logs"`
	} `json:"Info"`
}

// Получить данные по методу SearchStartRequest
func (user User) SearchStartRequest(searchStartReq SearchStartRequestStruct) (SearchStartResponseStruct, error) {
	searchStartReq.UserId = user.UserId
	searchStartReq.UserLogin = user.UserLogin
	searchStartReq.UserPassword = user.UserPassword

	// Ответ от сервера
	var responseSearchStart SearchStartResponseStruct

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(searchStartReq)
	if responseError != nil {
		return responseSearchStart, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "SearchStart")
	if responseError != nil {
		return responseSearchStart, responseError
	}

	// Распарсить данные
	responseError = responseSearchStart.UnmarshalJson(body)

	return responseSearchStart, responseError
}

// Метод для SearchStartResponseStruct, который преобразует приходящий ответ в структуру
func (responseSearchStart *SearchStartResponseStruct) UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseSearchStart)
	if responseError != nil {
		return responseError
	}

	if len(responseSearchStart.Info.Errors) != 0 {
		return errors.New(responseSearchStart.Info.Errors[0])
	}
	return nil
}
