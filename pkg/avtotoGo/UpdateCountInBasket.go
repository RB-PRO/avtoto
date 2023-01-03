package avtotoGo

import (
	"encoding/json"
	"strconv"
)

// Метод UpdateCountInBasket предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO. Расширенная версия, выдает статус ответа.

// Вся структура запроса метода UpdateCountInBasket
type updateCountInBasketRequestData struct {
	User  User                         `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []UpdateCountInBasketRequest `json:"parts"` // Список запчастей для обновления количества в корзине (тип: индексированный массив):
	// Примечание: Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Тело запроса UpdateCountInBasket
type UpdateCountInBasketRequest struct {
	InnerID  int  `json:"InnerID"`  // ID записи в корзине AvtoTO (тип: целое) — данные, сохраненные в результате добавления в корзину
	RemoteID int  `json:"RemoteID"` // ID запчасти в Вашей системе (тип: целое)
	NewCount uint `json:"NewCount"` // Необходимое количество запчастей в карзине (тип: целое, положительное) — Необходимо, чтобы новое количество NewCount не превышало максимальное количество MaxCount, и соответствовало кратности заказа BaseCount
	// Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Тело ответа AddToBasket
type UpdateCountInBasketResponse struct {
	Done []struct { // Массив RemoteID успешно добавленных элементов
		RemoteID int `json:"RemoteID"`
	} `json:"Done"`
	Errors []struct { // Массив ошибок:
		Type  string `json:"type"`  // Тип ошибки: RemoteID - Если элемент прошел проверку на корректность, но возникла ошибка при добавлении элемента в корзину или Element, если возникла ошибка при проверке на корректность
		Id    int    `json:"id"`    // RemoteID или номер элемента
		Error string `json:"error"` // Описание ошибки
	} `json:"Errors"`
	Info struct { // Общая информация по запросу
		DocVersion string `json:"DocVersion"` // Версия документации
		IP         string `json:"IP"`         // IP используемой машины
		UserID     int    `json:"UserID"`     // ID пользователя
	} `json:"Info"`
}

// Получить данные по методу UpdateCountInBasket
func (user User) UpdateCountInBasket(UpdateCountInBasketReq []UpdateCountInBasketRequest) (UpdateCountInBasketResponse, error) {
	UpdateCountInBasketData := updateCountInBasketRequestData{User: user, Parts: UpdateCountInBasketReq}

	// Ответ от сервера
	var UpdateCountInBasketRes UpdateCountInBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(UpdateCountInBasketData)
	if responseError != nil {
		return UpdateCountInBasketResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "UpdateCountInBasket")
	if responseError != nil {
		return UpdateCountInBasketResponse{}, responseError
	}

	// Распарсить данные
	responseError = UpdateCountInBasketRes.updateCountInBasket_UnmarshalJson(body)
	if responseError != nil {
		return UpdateCountInBasketResponse{}, responseError
	}
	return UpdateCountInBasketRes, nil
}

// Метод для UpdateCountInBasket, который преобразует приходящий ответ в структуру
func (UpdateCountInBasketRes *UpdateCountInBasketResponse) updateCountInBasket_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &UpdateCountInBasketRes)
	if responseError != nil {
		return responseError
	}
	return nil
}

// Получить ошибку из ответа метода UpdateCountInBasket
func (UpdateCountInBasketRes UpdateCountInBasketResponse) ErrorString() string {
	if len(UpdateCountInBasketRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range UpdateCountInBasketRes.Errors {
			exitString += "Тип ошибки " + valueBasketError.Type +
				", RemoteID  " + strconv.Itoa(valueBasketError.Id) +
				", Описание ошибки " + valueBasketError.Type + ". "
		}
		return exitString
	}
}
