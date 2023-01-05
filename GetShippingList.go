package avtotoGo

import (
	"encoding/json"
)

// Метод GetShippingList предназначен для получения статистики проценок по всем объединенным регистрациям.

// Тело запроса GetShippingList
type GetShippingListRequestData struct {
	User User `json:"user"` // Данные пользователя для авторизации (тип: ассоциативный массив)
	data GetShippingListRequest
}
type GetShippingListRequest struct {
	From    Date `json:"from,omitempty"`     // дата начала выборки (ДД.ММ.ГГГГ) (опционально) (тип: строка)
	To      Date `json:"to,omitempty"`       // дата окончания выборки (ДД.ММ.ГГГГ) (опционально) (тип: строка)
	PageNum int  `json:"page_num,omitempty"` // номер страницы (опционально) (тип: целое)
}

// Тело ответа GetShippingList
type GetShippingListResponse struct {
	Shippings []struct { // Список отгрузок - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
		Id     int        `json:"Id"`   // ID отгрузки
		Date   int        `json:"Date"` // Дата отгрузки
		Type   int        `json:"Type"` // Тип отгрузки
		Summ   int        `json:"Summ"` // Сумма отгрузки
		Orders []struct { // Список заказов отгрузки - индексированный массив с упорядоченными целочисленными ключами, начиная с 0
			OrderId int    `json:"OrderId"` // ID заказа
			Comment string `json:"Comment"` // Комментарий к заказу
			Sclad   string `json:"Sclad"`   // Направление склада
			Code    string `json:"Code"`    // Артикул
			Name    string `json:"Name"`    // Название
			Manuf   string `json:"Manuf"`   // Бренд
			Price   int    `json:"Price"`   // Цена
			Count   int    `json:"Count"`   // Количество
		} `json:"Orders"`
	} `json:"Shippings"`
	Pagination struct { // массив с информацией о пагинации
		CountItems  int `json:"CountItems"`  // Количество отгрузок
		CountPages  int `json:"CountPages"`  // Количество страниц
		PerPage     int `json:"PerPage"`     // Количество результатов на странице
		CurrentPage int `json:"CurrentPage"` // Текущая страница
	} `json:"Pagination"`
	Errors []string `json:"Errors"` // Массив ошибок, возникший в процессе поиска
	Info   struct { // Общая информация о запросе
		DocVersion string `json:"DocVersion"` // Версия документации
	} `json:"Info"`
}

// Получить данные по методу GetShippingList
func (user User) GetShippingList(GetShippingListReq GetShippingListRequest) (GetShippingListResponse, error) {
	GetShippingListData := GetShippingListRequestData{User: user, data: GetShippingListReq}

	// Ответ от сервера
	var GetShippingListRes GetShippingListResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(GetShippingListData)
	if responseError != nil {
		return GetShippingListResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "GetShippingList")
	if responseError != nil {
		return GetShippingListResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &GetShippingListRes)
	if responseErrorUnmarshal != nil {
		return GetShippingListResponse{}, responseErrorUnmarshal
	}

	return GetShippingListRes, nil
}

// Получить Количество отгрузок из ответа метода GetShippingList
func (slr GetShippingListResponse) CountItems() int {
	return slr.Pagination.CountItems
}

// Получить Количество страниц из ответа метода GetShippingList
func (slr GetShippingListResponse) CountPages() int {
	return slr.Pagination.CountPages
}

// Получить Количество результатов на странице из ответа метода GetShippingList
func (slr GetShippingListResponse) PerPage() int {
	return slr.Pagination.PerPage
}

// Получить Текущую страницу из ответа метода GetShippingList
func (slr GetShippingListResponse) CurrentPage() int {
	return slr.Pagination.CurrentPage
}

// Получить ошибки из ответа метода GetShippingList
func (slr GetShippingListResponse) Error() string {
	return slr.Errors[0]
}
