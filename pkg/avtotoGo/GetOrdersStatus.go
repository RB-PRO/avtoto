package avtotoGo

// Метод GetOrdersStatus проверка статуса заказа в системе AvtoTO

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Вся структура запроса метода GetOrdersStatus
type GetOrdersStatusRequestData struct {
	User  User                     `json:"user"`  // Данные пользователя для авторизации (тип: ассоциативный массив)
	Parts []GetOrdersStatusRequest `json:"parts"` // Список запчастей для добавления в заказы (тип: индексированный массив)
}

// Тело запроса GetOrdersStatus
type GetOrdersStatusRequest struct {
	InnerID  int `json:"InnerID"`  // ID записи в корзине AvtoTO (тип: целое) — данные, сохраненные в результате добавления в корзину
	RemoteID int `json:"RemoteID"` // ID запчасти в Вашей системе (тип: целое)
}

// Тело ответа GetOrdersStatus
type GetOrdersStatusResponse struct {
	OrdersInfo []struct { // Массив с информацией о статусах заказов
		RemoteID int      `json:"RemoteID"` // ID заказа в Вашей системе
		InnerID  int      `json:"InnerID"`  // ID заказа в системе AvtoTO
		Info     struct { // массив данных о статусе заказа
			Progress      int      `json:"progress"`      // общий статус заказа (тип: целое)
			Progress_text string   `json:"progress_text"` // общий статус заказа (тип: строка)
			Count         int      `json:"count"`         // общее количество заказа (тип: целое)
			Sub_progress  []string `json:"sub_progress"`  // частичные статусы заказа (тип: массив)
			// Частичный статус (номер) => количество частичного статуса
			Sub_progress_text string `json:"sub_progress_text"` // частичные статусы заказа, описание (тип: строка с HTML разметкой)
		} `json:"Info"`
	} `json:"OrdersInfo"`
	Errors []struct { // Массив ошибок:
		RemoteID int      `json:"RemoteID"` // ID товара в Вашей системе
		InnerID  int      `json:"InnerID"`  // ID товара в корзине AvtoTO (тип: целое)
		Errors   []string `json:"Errors"`   // список ошибок по данному ID товара (массив)
	} `json:"Errors"`
}

// Получить данные по методу GetOrdersStatus
func (user User) GetOrdersStatus(GetOrdersStatusReq []GetOrdersStatusRequest) (GetOrdersStatusResponse, error) {
	GetOrdersStatusData := GetOrdersStatusRequestData{User: user, Parts: GetOrdersStatusReq}

	// Ответ от сервера
	var GetOrdersStatusRes GetOrdersStatusResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(GetOrdersStatusData)
	if responseError != nil {
		return GetOrdersStatusResponse{}, responseError
	}

	// Отправить данные
	body, responseError := HttpPost(bytesRepresentation, "GetOrdersStatus")
	if responseError != nil {
		return GetOrdersStatusResponse{}, responseError
	}
	fmt.Println(string(body))
	// Распарсить данные
	responseError = GetOrdersStatusRes.GetOrdersStatus_UnmarshalJson(body)
	if responseError != nil {
		return GetOrdersStatusResponse{}, responseError
	}
	return GetOrdersStatusRes, nil
}

// Метод для GetOrdersStatus, который преобразует приходящий ответ в структуру
func (GetOrdersStatusRes *GetOrdersStatusResponse) GetOrdersStatus_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &GetOrdersStatusRes)
	if responseError != nil {
		return responseError
	}
	return nil
}

// Получить ошибку из ответа метода GetOrdersStatus
func (GetOrdersStatusRes GetOrdersStatusResponse) ErrorString() string {
	if len(GetOrdersStatusRes.Errors) == 0 {
		return ""
	} else {
		var exitString string
		for _, valueBasketError := range GetOrdersStatusRes.Errors {
			exitString += "ID свой " + strconv.Itoa(valueBasketError.RemoteID) +
				", ID корзины " + strconv.Itoa(valueBasketError.InnerID) +
				", ошибки " + strings.Join(valueBasketError.Errors, ";") + ". "
		}
		return exitString
	}
}

// Получить статус заказа
func (GetOrdersStatusRes GetOrdersStatusResponse) Status(partCount int) (string, error) {
	if len(GetOrdersStatusRes.OrdersInfo) >= partCount {
		switch GetOrdersStatusRes.OrdersInfo[partCount].Info.Progress {
		case 1:
			return "Ожидает обработки", nil
		case 2:
			return "Ожидает оплаты", nil
		case 3:
			return "Заказано", nil
		case 4:
			return "Закуплено", nil
		case 5:
			return "В пути", nil
		case 6:
			return "На складе", nil
		case 7:
			return "Выдано", nil
		case 8:
			return "Нет в наличии", nil
		default:
			return "null", errors.New("enother status GetOrdersStatusRes.OrdersInfo[partCount]")
		}
	} else {
		return "null", errors.New("enother len(GetOrdersStatusRes.OrdersInfo[partCount])")
	}
}
