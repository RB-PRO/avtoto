package avtotoGo

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Метод [AddToBasket] добавляет запчасти в корзину
// Полная структура запроса метода AddToBasket скрыта от разработчика.
//
// [AddToBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#AddToBasket
type addToBasketRequestData struct {
	User  User                 `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []AddToBasketRequest `json:"parts"` // Список запчастей для добавления в корзину (тип: индексированный массив)
}

// Метод [AddToBasket] добавляет запчасти в корзину:
//
// Примечание: Необходимо, чтобы количество для покупки Count не превышало максимальное количество MaxCount и соответствовало кратности заказа BaseCount. [*] — данные, сохраненные в результате поиска
//
// # Структура запроса метода AddToBasket
//
// [AddToBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#AddToBasket
type AddToBasketRequest struct {
	Code     string  `json:"Code"`              // [*] Код детали
	Manuf    string  `json:"Manuf"`             // [*] Производитель
	Name     string  `json:"Name"`              // [*] Название
	Price    float64 `json:"Price"`             // Цена
	Storage  string  `json:"Storage"`           // [*] Склад
	Delivery string  `json:"Delivery"`          // [*] Срок доставки
	Count    int     `json:"Count"`             // [*] количество для покупки (тип: целое)
	PartId   int     `json:"PartId"`            // [*] Номер запчасти в списке результата поиска (тип: целое)
	SearchID int     `json:"SearchID"`          // [*] Номер поиска (тип: целое)
	RemoteID int     `json:"RemoteID"`          // ID запчасти в Вашей системе(тип: целое)
	Comment  string  `json:"Comment,omitempty"` // Ваш комментарий к запчасти (тип: строка) [необязательный параметр]
}

// Метод [AddToBasket] добавляет запчасти в корзину:
//
// # Структура ответа метода AddToBasket
//
// [AddToBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#AddToBasket
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

// Функция преобразования ответа результата добавления товара в корзину в запрос на обновление позиции товара.
// Для примера была взята структура AddToBasketRes:
//
//	basketItemsUpdate, errorBasketItemUpdate := AddToBasketRes.BasketResInUpdateReq(0)
//	if errorBasketItemUpdate != nil {
//		fmt.Println(errorBasketItemUpdate)
//	}
//
// Output:
//
//	avtotoGo.UpdateCountInBasketRequest{InnerID:99756690, RemoteID:1, NewCount:0x0}
func (AddToBasketRes AddToBasketResponse) BasketResInUpdateReq(partCount int) (UpdateCountInBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return UpdateCountInBasketRequest{}, errors.New("length AddToBasketRes.DoneInnerID = 0")
	}
	return UpdateCountInBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Функция преобразования ответа результата добавления товара в корзину в запрос на удаление позиции товара.
// Для примера была взята структура AddToBasketRes:
//
//	basketItemsDelete, errorBasketItemDelete := AddToBasketRes.BasketResInDeleteReq(0)
//	if errorBasketItemDelete != nil {
//		fmt.Println(errorBasketItemDelete)
//	}
//
// Output:
//
//	avtotoGo.UpdateCountInBasketRequest{InnerID:99756690, RemoteID:1, NewCount:0x0}
func (AddToBasketRes AddToBasketResponse) BasketResInDeleteReq(partCount int) (DeleteFromBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return DeleteFromBasketRequest{}, errors.New("length AddToBasketRes.DoneInnerID = 0")
	}
	return DeleteFromBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на получение информации по товару из корзины
//
//	basketCheck, errorbasketChecks := AddToBasketRes.BasketResInCheckReq(0)
//	if errorbasketChecks != nil {
//		fmt.Println(errorbasketChecks)
//	}
//	fmt.Printf("%+#v\n", basketCheck)
//
// Output:
//
//	avtotoGo.CheckAvailabilityInBasketRequest{InnerID:99756690, RemoteID:1, Count:0}
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
//
//	orderBasket, errorbasketChecks := AddToBasketRes.BasketResInOrdersReq(0)
//	if errorbasketChecks != nil {
//		fmt.Println(errorbasketChecks)
//	}
//	fmt.Printf("%+#v\n", orderBasket)
//
// Output:
//
//	avtotoGo.AddToOrdersFromBasketRequest{InnerID:99756690, RemoteID:1, Count:0}
func (AddToBasketRes AddToBasketResponse) BasketResInOrdersReq(partCount int) (AddToOrdersFromBasketRequest, error) {
	if len(AddToBasketRes.DoneInnerID) == 0 {
		return AddToOrdersFromBasketRequest{}, errors.New("length AddToOrdersFromBasketRequest.DoneInnerID = 0")
	}
	return AddToOrdersFromBasketRequest{
		InnerID:  AddToBasketRes.DoneInnerID[partCount].InnerID,
		RemoteID: AddToBasketRes.DoneInnerID[partCount].RemoteID,
	}, nil
}

// Преобразовать ответ после добавления товара в корзину в запрос на получения статуса заказа
//
//	orderStatusGet, errorbasketChecks := AddToBasketRes.BasketResInOrdersStatusReq(0)
//	if errorbasketChecks != nil {
//		fmt.Println(errorbasketChecks)
//	}
//	fmt.Printf("%+#v\n", orderStatusGet)
//
// Output:
//
//	avtotoGo.GetOrdersStatusRequest{InnerID:99756690, RemoteID:1}
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
//
//	basketItems := make([]avtoto.AddToBasketRequest, 1)
//	basketItem, errorBasketItem := SearchGetParts2Res.SearchResInBasketReq(0)
//	if errorBasketItem != nil {
//		fmt.Println(errorBasketItem)
//	}
//	basketItems[0] = basketItem
//	basketItems[0].RemoteID = 1
//	basketItems[0].Count = 20
//	AddToBasketRes, errorRes := user.AddToBasket(basketItems)
//	if errorRes != nil {
//		fmt.Println(errorRes)
//	}
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
	body, responseError := httpPost(bytesRepresentation, "AddToBasket")
	if responseError != nil {
		return AddToBasketResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &AddToBasketRes)
	if responseErrorUnmarshal != nil {
		return AddToBasketResponse{}, responseErrorUnmarshal
	}

	return AddToBasketRes, nil
}

// Получить ошибку из ответа метода AddToBasket
func (AddToBasketRes AddToBasketResponse) Error() string {
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
