package avtotoGo

import (
	"encoding/json"
)

// Метод [GetBrandsByCode] предназначен для поиска списка брендов по артикулу запчасти
// Примечание: Сервис поиска предложений будет работать в случае выполнения условия: сумма заказов / количество запросов > 20 после некоторого порога проценок. [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
//
// # Структура запроса метода GetBrandsByCode
//
// [GetBrandsByCode]: https://www.avtoto.ru/services/search/docs/technical_soap.html#GetBrandsByCode
type GetBrandsByCodeRequest struct {
	UserId       int    `json:"user_id"`       // [*] Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // [*] Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
	SearchCode   string `json:"search_code"`   // Поисковый запрос, минимум 3 символа (тип: строка)
}

// Метод [GetBrandsByCode] предназначен для поиска списка брендов по артикулу запчасти
//
// # Структура ответа метода GetBrandsByCode
//
// [GetBrandsByCode]: https://www.avtoto.ru/services/search/docs/technical_soap.html#GetBrandsByCode
type GetBrandsByCodeResponse struct {
	Brands []struct { // Список брендов, найденных по запросу - индексированный массив с упорядоченными целочисленными ключами, начиная с 0.
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

// Получить данные по методу GetBrandsByCode
//
//	myBrand := avtoto.GetBrandsByCodeRequest{SearchCode: mySearchCode} // Создаём структуру запроса бренда по заданному артиклу
//	dataGetBrandsByCode, errorSearch := user.GetBrandsByCode(myBrand) // Получаем с сервера список брендов
//	if errorSearch != nil {
//		log.Fatal(errorSearch)
//	}
//	fmt.Println("> Для артикла", mySearchCode, "найдено", len(dataGetBrandsByCode.Brands), "бренда(ов).",
//		"\nПервый найденный бренд имеет производителя", dataGetBrandsByCode.Brands[0].Manuf, "и имя", dataGetBrandsByCode.Brands[0].Name)
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
	body, responseError := httpPost(bytesRepresentation, "GetBrandsByCode")
	if responseError != nil {
		return GetBrandsByCodeResponse{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &GetBrandsByCodeRes)
	if responseErrorUnmarshal != nil {
		return GetBrandsByCodeResponse{}, responseErrorUnmarshal
	}

	return GetBrandsByCodeRes, responseError
}

// Получить ошибку из ответа метода GetBrandsByCode
func (GetBrandsByCodeRes GetBrandsByCodeResponse) Error() string {
	if len(GetBrandsByCodeRes.Info.Errors) == 0 {
		return ""
	} else {
		return GetBrandsByCodeRes.Info.Errors[0]
	}
}
