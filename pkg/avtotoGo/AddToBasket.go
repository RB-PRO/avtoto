package avtotoGo

// Метод AddToBasket добавляет запчасти в корзину

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Вся структура запроса метода AddToBasket
type addToBasketRequestData struct {
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

	Count    int    `json:"Count"`             // [*] количество для покупки (тип: целое)
	PartId   int    `json:"PartId"`            // [*] Номер запчасти в списке результата поиска (тип: целое)
	SearchID int    `json:"SearchID"`          // [*] Номер поиска (тип: целое)
	RemoteID int    `json:"RemoteID"`          // ID запчасти в Вашей системе(тип: целое)
	Comment  string `json:"Comment,omitempty"` // Ваш комментарий к запчасти (тип: строка) [необязательный параметр]
	// [*] — данные, сохраненные в результате поиска
	// Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount
}

// Тело ответа AddToBasket
type AddToBasketResponse struct {
	Done   []int      `json:"Done"` // Массив RemoteID успешно добавленных элементов
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
	DoneInnerID []struct { // Массив успешно добавленных запчастей с внутренними ID корзины:
		RemoteID int `json:"RemoteID"` // ID товара в Вашей системе
		InnerID  int `json:"InnerID"`  // ID товара в корзине AvtoTO
	} `json:"DoneInnerId"`
}

// Преобразовать ответ после добавления товара в корзину в запрос на обновление
func (AddToBasketRes AddToBasketResponse) BasketResInUpdateReq(partCount int) (UpdateCountInBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return UpdateCountInBasketRequest{}, errors.New("length AddToBasketRes.DoneInnerID = 0")
	}
	return UpdateCountInBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на удаление
func (AddToBasketRes AddToBasketResponse) BasketResInDeleteReq(partCount int) (DeleteFromBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return DeleteFromBasketRequest{}, errors.New("length AddToBasketRes.DoneInnerID = 0")
	}
	return DeleteFromBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на получение информации
func (AddToBasketRes AddToBasketResponse) BasketResInCheckReq(partCount int) (CheckAvailabilityInBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return CheckAvailabilityInBasketRequest{}, errors.New("length AddToBasketRes.DoneInnerID = 0")
	}
	return CheckAvailabilityInBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на добавление запчасти из корзины в заказы
func (AddToBasketRes AddToBasketResponse) BasketResInOrdersReq(partCount int) (AddToOrdersFromBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return AddToOrdersFromBasketRequest{}, errors.New("length AddToOrdersFromBasketRequest.DoneInnerID = 0")
	}
	return AddToOrdersFromBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на добавление запчасти из корзины в заказы
func (AddToBasketRes AddToBasketResponse) BasketResInOrdersStatusReq(partCount int) (GetOrdersStatusRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return GetOrdersStatusRequest{}, errors.New("length GetOrdersStatusRequest.DoneInnerID = 0")
	}
	return GetOrdersStatusRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Получить данные по методу AddToBasket
func (user User) AddToBasket(AddToBasketReq []AddToBasketRequest) (AddToBasketResponse, error) {
	AddToBasketReqData := addToBasketRequestData{User: user, Parts: AddToBasketReq}

	// Ответ от сервера
	var AddToBasketRes AddToBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(AddToBasketReqData)
	if responseError != nil {
		return AddToBasketResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "AddToBasket")
	if responseError != nil {
		return AddToBasketResponse{}, responseError
	}

	//fmt.Println(string(body))

	//Распарсить данные
	responseError = AddToBasketRes.addToBasket_UnmarshalJson(body)
	if responseError != nil {
		return AddToBasketResponse{}, responseError
	}
	return AddToBasketRes, nil
}

// Метод для SearchGetParts2, который преобразует приходящий ответ в структуру
func (responseAddToBasket *AddToBasketResponse) addToBasket_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseAddToBasket)
	if responseError != nil {
		return responseError
	}

	//if len(responseAddToBasket.Info.Errors) != 0 {
	//	return errors.New(responseAddToBasket.Info.Errors[0])
	//}
	return nil
}

// Получить ошибку из ответа метода AddToBasket
func (AddToBasketRes AddToBasketResponse) ErrorString() string {
	if len(AddToBasketRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range AddToBasketRes.Errors {
			exitString += "Тип ошибки " + valueBasketError.Type +
				", RemoteID  " + strconv.Itoa(valueBasketError.Id) +
				", Описание ошибки " + valueBasketError.Type + ". "
		}
		return exitString
	}
}
