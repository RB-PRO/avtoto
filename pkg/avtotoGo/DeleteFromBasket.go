package avtotoGo

// Метод DeleteFromBasket удаляет запчасти из корзины

import (
	"encoding/json"
	"strconv"
)

// Вся структура запроса метода DeleteFromBasket
type deleteFromBasketRequestData struct {
	User  User                      `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []DeleteFromBasketRequest `json:"parts"` // Список запчастей для удаления из корзины (тип: индексированный массив):
}

// Тело запроса DeleteFromBasket
type DeleteFromBasketRequest struct {
	InnerID  int `json:"InnerID"`  // ID записи в корзине AvtoTO (тип: целое)
	RemoteID int `json:"RemoteID"` // ID запчасти в Вашей системе (тип: целое)
}

// Тело ответа DeleteFromBasket
type DeleteFromBasketResponse struct {
	//Done []struct { // Массив RemoteID успешно удаленных элементов
	//	RemoteID int `json:"RemoteID"`
	//} `json:"Done"`
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
	body, responseError := HttpPost(bytesRepresentation, "DeleteFromBasket")
	if responseError != nil {
		return DeleteFromBasketResponse{}, responseError
	}

	// Распарсить данные
	responseError = DeleteFromBasketRes.DeleteFromBasket_UnmarshalJson(body)
	if responseError != nil {
		return DeleteFromBasketResponse{}, responseError
	}
	return DeleteFromBasketRes, nil
}

// Метод для DeleteFromBasket, который преобразует приходящий ответ в структуру
func (DeleteFromBasketRes *DeleteFromBasketResponse) DeleteFromBasket_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &DeleteFromBasketRes)
	if responseError != nil {
		return responseError
	}
	return nil
}

// Получить ошибку из ответа метода DeleteFromBasket
func (DeleteFromBasketRes DeleteFromBasketResponse) ErrorString() string {
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
