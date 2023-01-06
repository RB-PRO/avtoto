package avtoto

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// Метод [CheckAvailabilityInBasket] проверяет запчасти в корзине AvtoTO на наличие в прайсах для дальнейшего заказа, а так же срок хранения в корзине
// Полная структура запроса метода CheckAvailabilityInBasket скрыта от разработчика.
//
// [CheckAvailabilityInBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#CheckAvailabilityInBasket
type checkAvailabilityInBasketRequestData struct {
	User  User                               `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []CheckAvailabilityInBasketRequest `json:"parts"` // Список запчастей для удаления из корзины (тип: индексированный массив):
}

// Метод [CheckAvailabilityInBasket] проверяет запчасти в корзине AvtoTO на наличие в прайсах для дальнейшего заказа, а так же срок хранения в корзине
//
// # Структура запроса метода CheckAvailabilityInBasket
//
// [CheckAvailabilityInBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#CheckAvailabilityInBasket
type CheckAvailabilityInBasketRequest struct {
	InnerID  int `json:"InnerID"`         // ID записи в корзине AvtoTO (тип: целое)
	RemoteID int `json:"RemoteID"`        // ID запчасти в Вашей системе (тип: целое)
	Count    int `json:"Count,omitempty"` // Количество для добавления (необязательный параметр, тип: целое)
}

// Метод [CheckAvailabilityInBasket] проверяет запчасти в корзине AvtoTO на наличие в прайсах для дальнейшего заказа, а так же срок хранения в корзине
//
// # Структура ответа метода CheckAvailabilityInBasket
//
// [CheckAvailabilityInBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#CheckAvailabilityInBasket
type CheckAvailabilityInBasketResponse struct {
	PartsInfo []struct { // информация о наличии товара в корзине, массив:
		RemoteID     int    `json:"RemoteID"`     // ID товара в Вашей системе
		InnerID      int    `json:"InnerID"`      // ID товара в корзине AvtoTO (тип: целое)
		Availability int    `json:"Availability"` // 1/0 (в наличии / нет в наличии) (тип: целое)
		MaxCount     string `json:"MaxCount"`     // максимальное допустимое количество товара для заказа в корзине AvtoTO (тип: целое, значение "-1" означает "без ограничений")
	} `json:"PartsInfo"`
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
}

// Получить данные по методу CheckAvailabilityInBasket
//
//	basketChecks := make([]avtoto.CheckAvailabilityInBasketRequest, 1)
//	basketCheck, errorbasketChecks := AddToBasketRes.BasketResInCheckReq(0)
//	if errorbasketChecks != nil {
//		fmt.Println(errorbasketChecks)
//	}
//	basketChecks[0] = basketCheck
//	CheckAvailabilityInBasketRes, errorCheckInBasket := user.CheckAvailabilityInBasket(basketChecks)
//	if errorCheckInBasket != nil {
//		fmt.Println(errorCheckInBasket)
//	}
//	availability, errorAvailability := CheckAvailabilityInBasketRes.Availability(0)
//	if errorAvailability != nil {
//		fmt.Println(errorAvailability)
//	}
func (user User) CheckAvailabilityInBasket(CheckAvailabilityInBasketReq []CheckAvailabilityInBasketRequest) (CheckAvailabilityInBasketResponse, error) {
	CheckAvailabilityInBasketData := checkAvailabilityInBasketRequestData{User: user, Parts: CheckAvailabilityInBasketReq}

	// Ответ от сервера
	var CheckAvailabilityInBasketRes CheckAvailabilityInBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(CheckAvailabilityInBasketData)
	if responseError != nil {
		return CheckAvailabilityInBasketResponse{}, responseError
	}

	// Отправить данные
	body, responseError := httpPost(bytesRepresentation, "CheckAvailabilityInBasket")
	if responseError != nil {
		return CheckAvailabilityInBasketResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &CheckAvailabilityInBasketRes)
	if responseErrorUnmarshal != nil {
		return CheckAvailabilityInBasketResponse{}, responseErrorUnmarshal
	}

	return CheckAvailabilityInBasketRes, nil
}

// Получить ошибку из ответа метода CheckAvailabilityInBasket
func (CheckAvailabilityInBasketRes CheckAvailabilityInBasketResponse) Error() string {
	if len(CheckAvailabilityInBasketRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range CheckAvailabilityInBasketRes.Errors {
			exitString += "ID свой " + strconv.Itoa(valueBasketError.RemoteID) +
				", ID корзины " + strconv.Itoa(valueBasketError.InnerID) +
				", ошибки " + strings.Join(valueBasketError.Errors, ";") + ". "
		}
		return exitString
	}
}

// Получить данные по Availability(наличие)
func (CheckAvailabilityInBasketRes CheckAvailabilityInBasketResponse) Availability(count int) (string, error) {
	if len(CheckAvailabilityInBasketRes.PartsInfo) >= count {
		if CheckAvailabilityInBasketRes.PartsInfo[count].Availability == 1 {
			return "в наличии", nil
		} else {
			return "нет в наличии", nil
		}
	} else {
		return "", errors.New("there is no given PartsInfo")
	}
}
