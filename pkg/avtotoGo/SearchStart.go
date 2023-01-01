package avtotoGo

// Метод SearchStart предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO. Расширенная версия, выдает статус ответа.

/*
Методы SearchStart и SearchGetParts позволяют организовать асинхронную передачу данных, помогая снизить нагрузку на Ваш и на наш сервер.

Метод SearchStart выдает идентификатор процесса поиска на сервере AvtoTO. Потом нужно отслеживать, не появился ли ответ с помощью метода SearchGetParts.
Эта схема работы реализуется на Ajax: первый запрос запускает метод SearchStart, по его окончанию вызывается функция, которая с небольшим периодом (0.3 - 0.5 сек) проверяет наличие ответа.
Когда ответ появился, она его выдает.
*/

import (
	"encoding/json"
	"errors"
)

// Структура запроса метода SearchStart
type SearchStartRequest struct {
	UserId       int    `json:"user_id"`         // [*] Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`      // [*] Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"`   // Пароль пользователя (тип: строка)
	SearchCode   string `json:"search_code"`     // Поисковый запрос, минимум 3 символа (тип: строка)
	SearchCross  string `json:"search_cross"`    // Искать в аналогах или нет (тип: строка, 'on' или 'off')
	Brand        string `json:"brand,omitempty"` // [**] Искать код с учетом бренда, минимум 2 символа (опционально)(тип: строка)
	// [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
	// [**] Список брендов можно получить с помощью метода GetBrandsByCode
	// Примечание: если бренд не указан, будет автоматически выбран самый популярный и произведен поиск с учетом этого бренда.
}

// Структура ответа метода SearchStart
type SearchStartResponse struct {
	ProcessSearchID string `json:"ProcessSearchId"` // идентификатор процесса поиска (тип: строка). Необходим для отслеживания результатов процесса поиска.
	Info            struct {
		SearchID string   `json:"SearchId"`
		Errors   []string `json:"Errors"`
		Logs     string   `json:"Logs"`
	} `json:"Info"`
}

// Преобразование ответа в запрос. SearchStartResponse > SearchGetParts2Request
func (searchStartRes SearchStartResponse) SearchResInReq() (SearchGetParts2Request, error) {
	// Проверка на Error в структуре ответа
	if len(searchStartRes.Info.Errors) != 0 {
		return SearchGetParts2Request{}, errors.New(searchStartRes.Info.Errors[0])
	}

	// Проверка существование ProcessSearchID
	if searchStartRes.ProcessSearchID == "" {
		return SearchGetParts2Request{}, errors.New("ProcessSearchID is nil")
	}

	return SearchGetParts2Request{ProcessSearchId: searchStartRes.ProcessSearchID}, nil
}

/* -------------------------------------------- */
/* ----**** JSON/http method functions ****---- */
/* -------------------------------------------- */

// Получить данные по методу SearchStartRequest
func (user User) SearchStartRequest(searchStartReq SearchStartRequest) (SearchStartResponse, error) {
	searchStartReq.UserId = user.UserId
	searchStartReq.UserLogin = user.UserLogin
	searchStartReq.UserPassword = user.UserPassword

	// Ответ от сервера
	var responseSearchStart SearchStartResponse

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
	responseError = responseSearchStart.searchStart_UnmarshalJson(body)

	return responseSearchStart, responseError
}

// Метод для SearchStartResponse, который преобразует приходящий ответ в структуру
func (responseSearchStart *SearchStartResponse) searchStart_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseSearchStart)
	if responseError != nil {
		return responseError
	}

	if len(responseSearchStart.Info.Errors) != 0 {
		return errors.New(responseSearchStart.Info.Errors[0])
	}
	return nil
}
