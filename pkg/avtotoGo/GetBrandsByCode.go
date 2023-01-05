package avtotoGo

import (
	"encoding/json"
)

// Метод GetBrandsByCode предназначен для поиска списка брендов по артикулу запчасти
// Примечание: Сервис поиска предложений будет работать в случае выполнения условия: сумма заказов / количество запросов > 20 после некоторого порога проценок.

// Структура запроса метода GetBrandsByCode
type GetBrandsByCodeRequest struct {
	UserId       int    `json:"user_id"`       // [*] Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // [*] Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
	SearchCode   string `json:"search_code"`   // Поисковый запрос, минимум 3 символа (тип: строка)
	// [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
}

// Структура ответа метода GetBrandsByCode
type GetBrandsByCodeResponse struct {
	// Список брендов, найденных по запросу - индексированный массив с упорядоченными целочисленными ключами, начиная с 0.
	// Каждый элемент этого массива содержит информацию о конкретном производителе и представляет из себя ассоциативный массив.
	// Свойства бренда:
	Brands []struct {
		Manuf string `json:"Manuf"` // Производитель
		Name  string `json:"Name"`  // Название
	} `json:"Brands"`
	Info struct {
		Errors []string `json:"Errors"`
	} `json:"Info"`
}

// Получить количество Parts метода GetBrandsByCode
func (GetBrandsByCodeRes GetBrandsByCodeResponse) LenParts() int {
	return len(GetBrandsByCodeRes.Brands)
}

/* -------------------------------------------- */
/* ----**** JSON/http method functions ****---- */
/* -------------------------------------------- */

// Получить данные по методу GetBrandsByCode
func (user User) GetBrandsByCode(GetBrandsByCodeReq GetBrandsByCodeRequest) (GetBrandsByCodeResponse, error) {
	GetBrandsByCodeReq.UserId = user.UserId
	GetBrandsByCodeReq.UserLogin = user.UserLogin
	GetBrandsByCodeReq.UserPassword = user.UserPassword

	// Ответ от сервера
	var GetBrandsByCodeRes GetBrandsByCodeResponse

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(GetBrandsByCodeReq)
	if responseError != nil {
		return GetBrandsByCodeResponse{}, responseError
	}

	// Выполнить запрос
	body, responseError := HttpPost(bytesRepresentation, "GetBrandsByCode")
	if responseError != nil {
		return GetBrandsByCodeResponse{}, responseError
	}

	// Распарсить данные
	responseError = GetBrandsByCodeRes.UnmarshalJSON(body)
	if responseError != nil {
		return GetBrandsByCodeResponse{}, responseError
	}

	return GetBrandsByCodeRes, responseError
}

// Метод для SearchStartResponse, который преобразует приходящий ответ в структуру

func (responseGetBrandsByCode *GetBrandsByCodeResponse) UnmarshalJSON(body []byte) error {
	//func (responseGetBrandsByCode *GetBrandsByCodeResponse) GetBrandsByCode_UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseGetBrandsByCode)
	if responseError != nil {
		return responseError
	}

	//if len(responseGetBrandsByCode.Info.Errors) != 0 {
	//	return errors.New(responseGetBrandsByCode.Info.Errors[0])
	//}
	return nil
}

// Получить ошибку из ответа метода GetBrandsByCode
func (GetBrandsByCodeRes GetBrandsByCodeResponse) ErrorString() string {
	if len(GetBrandsByCodeRes.Info.Errors) == 0 {
		return ""
	} else {
		return GetBrandsByCodeRes.Info.Errors[0]
	}
}
