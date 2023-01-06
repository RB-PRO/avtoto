package avtotoGo

import (
	"encoding/json"
	"strconv"
)

// Метод [DeleteFromBasket] удаляет запчасти из корзины
// Полная структура запроса метода DeleteFromBasket скрыта от разработчика.
//
// [DeleteFromBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#DeleteFromBasket
type deleteFromBasketRequestData struct {
	User  User                      `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []DeleteFromBasketRequest `json:"parts"` // Список запчастей для удаления из корзины (тип: индексированный массив):
}

// Метод [DeleteFromBasket] удаляет запчасти из корзины
//
// # Структура запроса метода DeleteFromBasket
//
// [DeleteFromBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#DeleteFromBasket
type DeleteFromBasketRequest struct {
	InnerID  int `json:"InnerID"`  // ID записи в корзине AvtoTO (тип: целое)
	RemoteID int `json:"RemoteID"` // ID запчасти в Вашей системе (тип: целое)
}

// Метод [DeleteFromBasket] удаляет запчасти из корзины
//
// # Структура ответа метода DeleteFromBasket
//
// [DeleteFromBasket]: https://www.avtoto.ru/services/search/docs/technical_soap.html#DeleteFromBasket
type DeleteFromBasketResponse struct {
	Done   []int      `json:"Done"` // Массив RemoteID успешно удаленных элементов
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

// Получить данные по методу DeleteFromBasket
//
//	basketItemsDeletes := make([]avtoto.DeleteFromBasketRequest, 1)
//	basketItemsDelete, errorBasketItemDelete := AddToBasketRes.BasketResInDeleteReq(0)
//	if errorBasketItemDelete != nil {
//		fmt.Println(errorBasketItemDelete)
//	}
//	basketItemsDeletes[0] = basketItemsDelete
//	DeleteFromBasketRes, errorBusketDelete := user.DeleteFromBasket(basketItemsDeletes)
//	if errorBasketItemDelete != nil {
//		fmt.Println(errorBusketDelete)
//	}
func (user User) DeleteFromBasket(DeleteFromBasketReq []DeleteFromBasketRequest) (DeleteFromBasketResponse, error) {
	DeleteFromBasketData := deleteFromBasketRequestData{User: user, Parts: DeleteFromBasketReq}

	// Ответ от сервера
	var DeleteFromBasketRes DeleteFromBasketResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(DeleteFromBasketData)
	if responseError != nil {
		return DeleteFromBasketResponse{}, responseError
	}

	// Отправить данные
	body, responseError := httpPost(bytesRepresentation, "DeleteFromBasket")
	if responseError != nil {
		return DeleteFromBasketResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &DeleteFromBasketRes)
	if responseErrorUnmarshal != nil {
		return DeleteFromBasketResponse{}, responseErrorUnmarshal
	}

	return DeleteFromBasketRes, nil
}

// Получить ошибку из ответа метода DeleteFromBasket
func (DeleteFromBasketRes DeleteFromBasketResponse) Error() string {
	if len(DeleteFromBasketRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range DeleteFromBasketRes.Errors {
			exitString += "Тип ошибки " + valueBasketError.Type +
				", RemoteID  " + strconv.Itoa(valueBasketError.Id) +
				", Описание ошибки " + valueBasketError.Type + ". "
		}
		return exitString
	}
}
