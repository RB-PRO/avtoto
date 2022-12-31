package avtotoGo

// Метод AddToBasket добавляет запчасти в корзину

import (
	"encoding/json"
	"fmt"
)

// Вся структура запроса метода AddToBasket
type AddToBasketRequestData struct {
	User  User                 `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []AddToBasketRequest `json:"parts"` // Список запчастей для добавления в корзину (тип: индексированный массив)
	// Примечание: Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Тело запроса AddToBasket
type AddToBasketRequest struct {
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
}

// Тело ответа AddToBasket
type AddToBasketResponse struct {
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
}

// Получить данные по методу AddToBasket
func (user User) AddToBasket(AddToBasketReq []AddToBasketRequest) (string, error) {
	AddToBasketReqData := AddToBasketRequestData{User: user, Parts: AddToBasketReq}

	// Ответ от сервера
	//var AddToBasketRes AddToBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(AddToBasketReqData)
	if responseError != nil {
		return "", responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "AddToBasket")
	if responseError != nil {
		return "", responseError
	}

	fmt.Println(string(body))

	return string(body), nil
	// Распарсить данные
	//responseError = AddToBasketRes.AddToBasket_UnmarshalJson(body)

	//return AddToBasketRes, responseError
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
