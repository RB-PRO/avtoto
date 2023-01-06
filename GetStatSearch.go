package avtoto

import (
	"encoding/json"
)

// Метод [GetStatSearch] предназначен для получения статистики проценок по всем объединенным регистрациям.
// Информаци о проценках - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
//
// # Структура ответа метода GetStatSearch
//
// [GetStatSearch]: https://www.avtoto.ru/services/search/docs/technical_soap.html#GetStatSearch
type GetStatSearchResponse struct {
	StatInfo struct {
		SearchCount   int  `json:"SearchCount"`   // Количество проценок за определенный период
		SearchEnabled bool `json:"SearchEnabled"` // Доступность использования проценки (true - доступно, false - недоступно)
		MaxCount      bool `json:"MaxCount"`      // лимит проценок
		OrdersSum     int  `json:"OrdersSum"`     // сумма закупок за определенный период

		StatDateStart      Date     `json:"StatDateStart"`      // дата начала периода подсчета
		StatDateStartStamp TimeUnix `json:"StatDateStartStamp"` // дата начала периода подсчета в формате UNIX

		StatDateEnd      Date     `json:"StatDateEnd"`      // Дата окончания периода подсчета
		StatDateEndStamp TimeUnix `json:"StatDateEndStamp"` // Дата окончания периода подсчета в формате UNIX

		SearchHistory []struct { // Информация о количестве проценок по дням - Массив со след. элементами:
			Day         Date `json:"Day"`         // День (в формате dd/mm)
			SearchCount int  `json:"SearchCount"` // Количество проценок
		} `json:"SearchHistory"`
	} `json:"StatInfo"`
	BrandsStatInfo struct { // Информаци о запросах брендов по коду - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
		SearchCount   string `json:"SearchCount"`   // Количество запросов за определенный период
		SearchEnabled bool   `json:"SearchEnabled"` // Доступность использования запросов (true - доступно, false - недоступно)
		MaxCount      bool   `json:"MaxCount"`      // Лимит запросов

		StatDateStart      Date       `json:"StatDateStart"`      // Дата начала периода подсчета
		StatDateStartStamp TimeUnix   `json:"StatDateStartStamp"` // Дата начала периода подсчета в формате UNIX
		StatDateEnd        Date       `json:"StatDateEnd"`        // Дата окончания периода подсчета
		StatDateEndStamp   TimeUnix   `json:"StatDateEndStamp"`   // Дата окончания периода подсчета в формате UNIX
		SearchHistory      []struct { // Информация о количестве запросов по дням - Массив со след. элементами:
			Day         Date `json:"Day"`         // День (в формате dd/mm)
			SearchCount int  `json:"SearchCount"` // Количество запросов
		} `json:"SearchHistory"`
	} `json:"BrandsStatInfo"`
	Errors []string `json:"Errors"` // Массив ошибок, возникший в процессе поиска
	Info   struct { // Общая информация по запросу
		DocVersion string `json:"DocVersion"` // Версия API
	} `json:"Info"`
}

// Получить данные по методу GetStatSearch
//
//	statSearch, statSearchError := user.GetStatSearch()
//	if statSearchError != nil {
//		fmt.Println(statSearchError)
//	}
func (user User) GetStatSearch() (GetStatSearchResponse, error) {

	// Ответ от сервера
	var GetStatSearchRes GetStatSearchResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(user)
	if responseError != nil {
		return GetStatSearchResponse{}, responseError
	}

	// Отправить данные
	body, responseError := httpPost(bytesRepresentation, "GetStatSearch")
	if responseError != nil {
		return GetStatSearchResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &GetStatSearchRes)
	if responseErrorUnmarshal != nil {
		return GetStatSearchResponse{}, responseErrorUnmarshal
	}

	return GetStatSearchRes, nil
}

// Получить ошибку из ответа метода GetStatSearch
func (GetStatSearchRes GetStatSearchResponse) Error() string {
	if len(GetStatSearchRes.Errors) == 0 {
		return ""
	} else {
		return GetStatSearchRes.Errors[0]
	}
}
