package avtoto

import (
	"encoding/json"
	"errors"
)

// Метод [SearchStart] предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO
// Расширенная версия, выдает статус ответа
// Метод SearchStart выдает идентификатор процесса поиска на сервере AvtoTO. Потом нужно отслеживать, не появился ли ответ с помощью метода SearchGetParts.
// Эта схема работы реализуется на Ajax: первый запрос запускает метод SearchStart, по его окончанию вызывается функция, которая с небольшим периодом (0.3 - 0.5 сек) проверяет наличие ответа.
// Когда ответ появился, она его выдает.
//
// # Структура запроса метода SearchStart
//
// [SearchStart]: https://www.avtoto.ru/services/search/docs/technical_soap.html#SearchStart

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

// Метод [SearchStart] предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO
//
// # Структура ответа метода SearchStart
//
// [SearchStart]: https://www.avtoto.ru/services/search/docs/technical_soap.html#SearchStart
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
	if searchStartRes.ProcessSearchID == "" { // Проверка существования ProcessSearchID
		return SearchGetParts2Request{}, errors.New("ProcessSearchID is nil")
	}
	return SearchGetParts2Request{ProcessSearchId: searchStartRes.ProcessSearchID}, nil
}

// Получить ProcessSearchID из ответа метода SearchStart
func (searchStartRes SearchStartResponse) ProcessSearchCode() string {
	return searchStartRes.ProcessSearchID
}

// Получить данные по методу SearchStart
//
//	searchStartReq := avtoto.SearchStartRequest{SearchCode: mySearchCode, SearchCross: "on", Brand: dataGetBrandsByCode.Brands[0].Manuf}// Объявление запроса метода SearchStart
//	datasSearchStartRequest, errorSearch := user.SearchStartRequest(searchStartReq)// Вызов метода SearchStartRequest с запросом
//	if errorSearch != nil {
//		log.Fatal(errorSearch)
//	}
//	fmt.Println("> Полученный ProcessSearchID", datasSearchStartRequest.ProcessSearchID)
func (user User) SearchStartRequest(searchStartReq SearchStartRequest) (SearchStartResponse, error) {
	searchStartReq.UserId = user.UserId
	searchStartReq.UserLogin = user.UserLogin
	searchStartReq.UserPassword = user.UserPassword

	// Ответ от сервера
	var SearchStartRes SearchStartResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(searchStartReq)
	if responseError != nil {
		return SearchStartResponse{}, responseError
	}

	// Выполнить запрос
	body, responseError := httpPost(bytesRepresentation, "SearchStart")
	if responseError != nil {
		return SearchStartResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &SearchStartRes)
	if responseErrorUnmarshal != nil {
		return SearchStartResponse{}, responseErrorUnmarshal
	}

	return SearchStartRes, responseError
}

// Получить ошибку из ответа метода SearchStart
func (SearchStartRes SearchStartResponse) Error() string {
	if len(SearchStartRes.Info.Errors) == 0 {
		return ""
	} else {
		return SearchStartRes.Info.Errors[0]
	}
}

// Получить логи из ответа метода SearchStart
func (responseSearchStart SearchStartResponse) LogsString() string {
	return responseSearchStart.Info.Logs
}
