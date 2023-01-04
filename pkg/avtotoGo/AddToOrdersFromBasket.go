package avtotoGo

// Метод AddToOrdersFromBasket добавляет запчасти в заказы из корзины Avtoto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Вся структура запроса метода AddToOrdersFromBasket
type AddToOrdersFromBasketRequestData struct {
	User  User                           `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []AddToOrdersFromBasketRequest `json:"parts"` // Список запчастей для удаления из корзины (тип: индексированный массив):
}

// Тело запроса AddToOrdersFromBasket
type AddToOrdersFromBasketRequest struct {
	InnerID  int `json:"InnerID"`         // ID записи в корзине AvtoTO (тип: целое) — данные, сохраненные в результате добавления в корзину
	RemoteID int `json:"RemoteID"`        // ID запчасти в Вашей системе (тип: целое)
	Count    int `json:"Count,omitempty"` // Количество для добавления (необязательный параметр, тип: целое)
}

// Тело ответа AddToOrdersFromBasket
type AddToOrdersFromBasketResponse struct {
	Done   []int      `json:"Done"` // Массив RemoteID успешно добавленных элементов
	Errors []struct { // Массив ошибок:
		RemoteID int      `json:"RemoteID"` // ID товара в Вашей системе
		InnerID  int      `json:"InnerID"`  // ID товара в корзине AvtoTO (тип: целое)
		Errors   []string `json:"Errors"`   // список ошибок по данному ID товара (массив)
	} `json:"Errors"`
	Info struct { // Общая информация по запросу
		DocVersion string `json:"DocVersion"` // Версия документации
		IP         string `json:"IP"`         // IP используемой машины
		UserID     int    `json:"UserID"`     // ID пользователя
	} `json:"Info"`
	DoneInnerId []struct { // Массив успешно добавленных запчастей с внутренними ID корзины:
		RemoteID int `json:"RemoteID"` // ID товара в Вашей системе
		InnerID  int `json:"InnerID"`  // InnerID - ID успешно добавленного в заказы товара AvtoTO
	} `json:"DoneInnerId"`
}

// Получить данные по методу AddToOrdersFromBasket
func (user User) AddToOrdersFromBasket(AddToOrdersFromBasketReq []AddToOrdersFromBasketRequest) (AddToOrdersFromBasketResponse, error) {
	AddToOrdersFromBasketData := AddToOrdersFromBasketRequestData{User: user, Parts: AddToOrdersFromBasketReq}

	// Ответ от сервера
	var AddToOrdersFromBasketRes AddToOrdersFromBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(AddToOrdersFromBasketData)
	if responseError != nil {
		return AddToOrdersFromBasketResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "AddToOrdersFromBasket")
	if responseError != nil {
		return AddToOrdersFromBasketResponse{}, responseError
	}
	fmt.Println(string(body))
	// Распарсить данные
	responseError = AddToOrdersFromBasketRes.AddToOrdersFromBasket_UnmarshalJson(body)
	if responseError != nil {
		return AddToOrdersFromBasketResponse{}, responseError
	}
	return AddToOrdersFromBasketRes, nil
}

// Метод для AddToOrdersFromBasket, который преобразует приходящий ответ в структуру
func (AddToOrdersFromBasketRes *AddToOrdersFromBasketResponse) AddToOrdersFromBasket_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &AddToOrdersFromBasketRes)
	if responseError != nil {
		return responseError
	}
	return nil
}

// Получить ошибку из ответа метода AddToOrdersFromBasket
func (AddToOrdersFromBasketRes AddToOrdersFromBasketResponse) ErrorString() string {
	if len(AddToOrdersFromBasketRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range AddToOrdersFromBasketRes.Errors {
			exitString += "ID свой " + strconv.Itoa(valueBasketError.RemoteID) +
				", ID корзины " + strconv.Itoa(valueBasketError.InnerID) +
				", ошибки " + strings.Join(valueBasketError.Errors, ";") + ". "
		}
		return exitString
	}
}
