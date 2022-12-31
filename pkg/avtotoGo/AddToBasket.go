package avtotoGo

import (
	"encoding/json"
	"fmt"
)

// Метод AddToBasket добавляет запчасти в корзину

// Структура запроса метода AddToBasket
type AddToBasketRequest struct {
	User  User       `json:"user"` // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []struct { // Список запчастей для добавления в корзину (тип: индексированный массив):
		Code     string  `json:"Code"`     // [*] Код детали
		Manuf    string  `json:"Manuf"`    // [*] Производитель
		Name     string  `json:"Name"`     // [*] Название
		Price    float64 `json:"Price"`    // Цена
		Storage  string  `json:"Storage"`  // [*] Склад
		Delivery string  `json:"Delivery"` // [*] Срок доставки

		Count    string `json:"Count"`    // [*] количество для покупки (тип: целое)
		PartId   string `json:"PartId"`   // [*] Номер запчасти в списке результата поиска (тип: целое)
		SearchID string `json:"SearchID"` // [*] Номер поиска (тип: целое)
		RemoteID string `json:"RemoteID"` // ID запчасти в Вашей системе(тип: целое)
		Comment  string `json:"Comment "` // Ваш комментарий к запчасти (тип: строка) [необязательный параметр]

		// [*] — данные, сохраненные в результате поиска
	} `json:"parts"`
	// Примечание: Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Структура ответа метода AddToBasket
type AddToBasketResponse struct {
	User  User       `json:"user"` // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []struct { // Список запчастей для добавления в корзину (тип: индексированный массив):
		Code     string  `json:"Code"`     // [*] Код детали
		Manuf    string  `json:"Manuf"`    // [*] Производитель
		Name     string  `json:"Name"`     // [*] Название
		Price    float64 `json:"Price"`    // Цена
		Storage  string  `json:"Storage"`  // [*] Склад
		Delivery string  `json:"Delivery"` // [*] Срок доставки

		Count    string `json:"Count"`    // [*] количество для покупки (тип: целое)
		PartId   string `json:"PartId"`   // [*] Номер запчасти в списке результата поиска (тип: целое)
		SearchID string `json:"SearchID"` // [*] Номер поиска (тип: целое)
		RemoteID string `json:"RemoteID"` // ID запчасти в Вашей системе(тип: целое)
		Comment  string `json:"Comment "` // Ваш комментарий к запчасти (тип: строка) [необязательный параметр]

		// [*] — данные, сохраненные в результате поиска
	} `json:"parts"`
	// Примечание: Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Получить данные по методу AddToBasket
func (user User) AddToBasket(AddToBasketReq AddToBasketRequest) (AddToBasketResponse, error) {
	AddToBasketReq.User.UserId = user.UserId
	AddToBasketReq.User.UserLogin = user.UserLogin
	AddToBasketReq.User.UserPassword = user.UserPassword

	// Ответ от сервера
	var AddToBasketRes AddToBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(AddToBasketReq)
	if responseError != nil {
		return AddToBasketRes, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "AddToBasket")
	if responseError != nil {
		return AddToBasketRes, responseError
	}

	fmt.Println(string(body))

	// Распарсить данные
	responseError = AddToBasketRes.AddToBasket_UnmarshalJson(body)

	return AddToBasketRes, responseError
}

// Метод для SearchGetParts2, который преобразует приходящий ответ в структуру
func (responseAddToBasket *AddToBasketResponse) AddToBasket_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseAddToBasket)
	if responseError != nil {
		return responseError
	}

	//if len(responseAddToBasket.Info.Errors) != 0 {
	//	return errors.New(responseAddToBasket.Info.Errors[0])
	//}
	return nil
}
