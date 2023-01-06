package avtotoGo

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Метод [SearchGetParts2] предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO. Расширенная версия, выдает статус ответа
// Полная структура запроса метода GetShippingList скрыта от разработчика.
//
// [SearchGetParts2]: https://www.avtoto.ru/services/search/docs/technical_soap.html#SearchGetParts2
type SearchGetParts2Request struct {
	ProcessSearchId string `json:"ProcessSearchId"` // Уникальный идентификатор процесса поиска (тип: строка).
	Limit           int    `json:"Limit"`           // необязательный параметр, орграничение на количество строк в выдаче (тип: целое).
}

// Метод [SearchGetParts2] предназначен для получения результатов поиска запчастей по коду на сервере AvtoTO. Расширенная версия, выдает статус ответа
// Список запчастей, найденных по запросу - индексированный массив с упорядоченными целочисленными ключами, начиная с 0.
// Каждый элемент этого массива содержит информацию о конкретной детали и представляет из себя ассоциативный массив.
// [*] — эти данные необходимо сохранить в Вашей системе, в дальнейшем они понадобятся для добавления запчастей в корзину
// [**] — В случае, когда SearchStatus = 4 (Результат получен)
//
// # Структура ответа метода SearchGetParts2
//
// [SearchGetParts2]: https://www.avtoto.ru/services/search/docs/technical_soap.html#SearchGetParts2
type SearchGetParts2Response struct {
	Parts []struct {
		Code      string `json:"Code"`      // [*] Код детали
		Manuf     string `json:"Manuf"`     // [*] Производитель
		Name      string `json:"Name"`      // [*] Название
		Price     int    `json:"Price"`     // Цена
		Storage   string `json:"Storage"`   // [*] Склад
		Delivery  string `json:"Delivery"`  // [*] Срок доставки
		MaxCount  string `json:"MaxCount"`  // [*] Максимальное количество для заказа, остаток по складу. Значение "-1" - означает "много" или "неизвестно"
		BaseCount string `json:"BaseCount"` // [*] Кратность заказа

		StorageDate     string `json:"StorageDate"`     // [**] Дата обновления склада
		DeliveryPercent int    `json:"DeliveryPercent"` // [**] Процент успешных закупок из общего числа заказов
		BackPercent     int    `json:"BackPercent"`     // [**] Процент удержания при возврате товара (при отсутствии возврата поставщику возвращается значение "-1")

		AvtotoData struct { // Массив со след. элементами:
			PartId int `json:"PartId"` // [*] Номер запчасти в списке результата поиска
		} `json:"AvtotoData"`
	} `json:"Parts"`
	Info struct {
		Errors       []string          `json:"Errors"`       // Массив ошибок, возникший в процессе поиска
		SearchStatus int               `json:"SearchStatus"` // Информация о статусе процесса на сервере AvtoTO. Возможные варианты значений:
		SearchID     CustomIntToString `json:"SearchId"`     // Уникальный идентификатор запроса поиска, возвращается в случае удачного поиска
	} `json:"Info"`
}

// Получить данные по методу SearchGetParts2
func (SearchGetParts2Res SearchGetParts2Response) SearchResInBasketReq(partCount int) (AddToBasketRequest, error) {
	SearchID_int, errorAtoi := SearchGetParts2Res.Info.SearchID.integer()
	if errorAtoi != nil {
		return AddToBasketRequest{}, errorAtoi
	}
	return AddToBasketRequest{
		Code: SearchGetParts2Res.Parts[partCount].Code,

		Manuf:    SearchGetParts2Res.Parts[partCount].Manuf,
		Name:     SearchGetParts2Res.Parts[partCount].Name,
		Price:    float64(SearchGetParts2Res.Parts[partCount].Price),
		Storage:  SearchGetParts2Res.Parts[partCount].Storage,
		Delivery: SearchGetParts2Res.Parts[partCount].Delivery,

		Count:    1,
		PartId:   SearchGetParts2Res.Parts[partCount].AvtotoData.PartId,
		SearchID: SearchID_int,
		//RemoteID: 1,
		//Comment: "",
	}, nil
}

// Структура созданная для десериализация JSON с неправильной типизацией - https://habr.com/ru/post/502176/
// Так получилось, что API можетотдавать данные по ключу SearchID как string, так и integer. Эта структура и 2 объявленных для неё метода способны изменить это
// и предоставить постоянный рабочий функционал, который позволяет держать значение элемента массива типа string в завосимости от входного параметра.
type CustomIntToString struct {
	value string
}

// Кастомное декодирование JSON для ключа SearchID
func (cis *CustomIntToString) UnmarshalJSON(data []byte) error {
	if data[0] == 34 { // Если первый символ - Кавычка
		err := json.Unmarshal(data, &cis.value)
		if err != nil {
			return errors.New("CustomIntToString: UnmarshalJSON: Find 34: " + err.Error())
		}
	} else {
		// Добавление Кавычек в начале и в конце массива byte
		newData := make([]byte, 1)
		newData[0] = 34
		newData = append(newData, data...)
		newData = append(newData, 34)
		//newData[len(data)] = 34

		err := json.Unmarshal(newData, &cis.value)
		if err != nil {
			return errors.New("CustomIntToString: UnmarshalJSON: Find't 34: " + err.Error())
		}
	}
	return nil
}

// Кастомное кодирование JSON для ключа SearchID
func (cf CustomIntToString) MarshalJSON() ([]byte, error) {
	json, err := json.Marshal(cf.value)
	return json, err
}

// Вернуть строку для SearchID. Необходимо для преобразования SearchGetParts2.Parts[int] > AddToBasket
func (cf CustomIntToString) integer() (int, error) {
	return strconv.Atoi(cf.value)
}

// Получить значение
func (cf CustomIntToString) Value() string {
	return cf.value
}

// Получить количество Parts метода SearchGetParts2
func (SearchGetParts2Res SearchGetParts2Response) LenParts() int {
	return len(SearchGetParts2Res.Parts)
}

// Получить данные по методу SearchGetParts2
//
//	SearchGetParts2Res, errorSearch := SearchGetParts2Req.SearchGetParts2() // Вызов метода SearchGetParts2
//	if errorSearch != nil {
//		log.Fatal(errorSearch)
//	}
//	if SearchGetParts2Res.Error() != "Запрос в обработке" {
//		fmt.Println(SearchGetParts2Res)
//	} else {
//		fmt.Println("Запрос в обработке. Ждём 1 секунду и заново опрашиваешь по методу SearchGetParts2")
//	}
func (SearchGetParts2Req SearchGetParts2Request) SearchGetParts2() (SearchGetParts2Response, error) {
	// Ответ от сервера
	var SearchGetParts2Res SearchGetParts2Response

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(SearchGetParts2Req)
	if responseError != nil {
		return SearchGetParts2Response{}, responseError
	}

	// Выполнить запрос
	body, responseError := httpPost(bytesRepresentation, "SearchGetParts2")
	if responseError != nil {
		return SearchGetParts2Response{}, responseError
	}

	// Распарсить данные
	responseErrorUnmarshal := json.Unmarshal(body, &SearchGetParts2Res)
	if responseErrorUnmarshal != nil {
		return SearchGetParts2Response{}, responseErrorUnmarshal
	}

	return SearchGetParts2Res, responseError
}

// Получить статус запроса по методу SearchGetParts2
func (SearchGetParts2Res SearchGetParts2Response) Status() string {
	switch SearchGetParts2Res.Info.SearchStatus {
	case 0:
		return "Неверно указан ID процесса ProcessSearchId"
	case 1:
		return "Запрос не найден"
	case 2:
		return "Запрос в обработке"
	case 3:
		return "Ошибка данных"
	case 4:
		return "Результат получен"
	default:
		return "Another error"
	}
}

// Получить ошибку из ответа метода SearchGetParts2
func (SearchGetParts2Res SearchGetParts2Response) Error() string {
	if len(SearchGetParts2Res.Info.Errors) == 0 {
		return ""
	} else {
		return SearchGetParts2Res.Info.Errors[0]
	}
}
