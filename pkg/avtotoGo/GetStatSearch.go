package avtotoGo

// Метод GetStatSearch предназначен для получения статистики проценок по всем объединенным регистрациям.

import (
	"encoding/json"
)

type GetStatSearchResponse struct {
	StatInfo struct { // Информаци о проценках - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
		SearchCount   int  `json:"SearchCount"`   // Количество проценок за определенный период
		SearchEnabled bool `json:"SearchEnabled"` // Доступность использования проценки (true - доступно, false - недоступно)
		MaxCount      bool `json:"MaxCount"`      // лимит проценок
		OrdersSum     int  `json:"OrdersSum"`     // сумма закупок за определенный период

		StatDateStart      string   `json:"StatDateStart"`      // дата начала периода подсчета
		StatDateStartStamp TimeUnix `json:"StatDateStartStamp"` // дата начала периода подсчета в формате UNIX

		StatDateEnd      string   `json:"StatDateEnd"`      // Дата окончания периода подсчета
		StatDateEndStamp TimeUnix `json:"StatDateEndStamp"` // Дата окончания периода подсчета в формате UNIX

		SearchHistory []struct { // Информация о количестве проценок по дням - Массив со след. элементами:
			Day         string `json:"Day"`         // День (в формате dd/mm)
			SearchCount int    `json:"SearchCount"` // Количество проценок
		} `json:"SearchHistory"`
	} `json:"StatInfo"`
	BrandsStatInfo struct { // Информаци о запросах брендов по коду - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
		SearchCount   string `json:"SearchCount"`   // Количество запросов за определенный период
		SearchEnabled bool   `json:"SearchEnabled"` // Доступность использования запросов (true - доступно, false - недоступно)
		MaxCount      bool   `json:"MaxCount"`      // Лимит запросов

		StatDateStart      string     `json:"StatDateStart"`      // Дата начала периода подсчета
		StatDateStartStamp int        `json:"StatDateStartStamp"` // Дата начала периода подсчета в формате UNIX
		StatDateEnd        string     `json:"StatDateEnd"`        // Дата окончания периода подсчета
		StatDateEndStamp   string     `json:"StatDateEndStamp"`   // Дата окончания периода подсчета в формате UNIX
		SearchHistory      []struct { // Информация о количестве запросов по дням - Массив со след. элементами:
			Day         string `json:"Day"`         // День (в формате dd/mm)
			SearchCount int    `json:"SearchCount"` // Количество запросов
		} `json:"SearchHistory"`
	} `json:"BrandsStatInfo"`
	Errors []string `json:"Errors"` // Массив ошибок, возникший в процессе поиска
	Info   struct { // Общая информация по запросу
		DocVersion string `json:"DocVersion"` // Версия API
	} `json:"Info"`
}

// Получить данные по методу GetStatSearch
func (user User) GetStatSearch() (GetStatSearchResponse, error) {

	// Ответ от сервера
	var GetStatSearchRes GetStatSearchResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(user)
	if responseError != nil {
		return GetStatSearchResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "GetStatSearch")
	if responseError != nil {
		return GetStatSearchResponse{}, responseError
	}

	// Распарсить данные
	responseError = GetStatSearchRes.GetStatSearch_UnmarshalJson(body)
	if responseError != nil {
		return GetStatSearchResponse{}, responseError
	}
	return GetStatSearchRes, nil
}

// Метод для GetStatSearch, который преобразует приходящий ответ в структуру
func (GetStatSearchRes *GetStatSearchResponse) GetStatSearch_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &GetStatSearchRes)
	if responseError != nil {
		return responseError
	}
	return nil
}
