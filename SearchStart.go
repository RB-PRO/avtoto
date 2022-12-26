package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

/*
	func encode_SearchStart(data []byte) ProcessSearchId {
		var result ProcessSearchId
		jsonErr := json.Unmarshal(data, &result)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return result
	}
*/

type SearchStartRequestStruct struct {
	UserId       int    `json:"user_id"`       // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
	SearchCode   string `json:"search_code"`   // Поисковый запрос, минимум 3 символа (тип: строка)
	SearchCross  string `json:"search_cross"`  // Искать в аналогах или нет (тип: строка, 'on' или 'off')
	Brand        string `json:"brand"`         /// Искать код с учетом бренда, минимум 2 символа (опционально)(тип: строка)
	// [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
	// [**] Список брендов можно получить с помощью метода GetBrandsByCode
}
type SearchStartResponseStruct struct {
	ProcessSearchID string `json:"ProcessSearchId"`
	Info            struct {
		SearchID string   `json:"SearchId"`
		Errors   []string `json:"Errors"`
		Logs     string   `json:"Logs"`
	} `json:"Info"`
}

func (user User) SearchStartRequest(searchStartReq SearchStartRequestStruct) (SearchStartResponseStruct, error) {
	searchStartReq.UserId = user.UserId
	searchStartReq.UserLogin = user.UserLogin
	searchStartReq.UserPassword = user.UserPassword

	// Ответ от сервера
	var responseSearchStart SearchStartResponseStruct

	// Подготовить данные для загрузки
	bytesRepresentation, responseError := json.Marshal(searchStartReq)
	if responseError != nil {
		return responseSearchStart, responseError
	}

	// Отправить данные
	body, responseError := httpPost(bytesRepresentation, "SearchStart")
	if responseError != nil {
		return responseSearchStart, responseError
	}

	// Распарсить данные
	responseError = responseSearchStart.UnmarshalJson(body)

	return responseSearchStart, responseError
}

func (responseSearchStart *SearchStartResponseStruct) UnmarshalJson(body []byte) error {
	responseError := json.Unmarshal(body, &responseSearchStart)
	if responseError != nil {
		return responseError
	}

	if len(responseSearchStart.Info.Errors) != 0 {
		return errors.New(responseSearchStart.Info.Errors[0])
	}
	return nil
}

func httpPost(bytesRepresentation []byte, action string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("action", "SearchStart")
	_ = writer.WriteField("data", string(bytesRepresentation))
	responseError := writer.Close()
	if responseError != nil {
		return nil, responseError
	}

	client := &http.Client{}
	req, responseError := http.NewRequest(http.MethodPost, URL, payload)
	if responseError != nil {
		return nil, responseError
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, responseError := client.Do(req)
	if responseError != nil {
		return nil, responseError
	}
	defer res.Body.Close()

	// Считываем ответ
	body, responseError := ioutil.ReadAll(res.Body)
	if responseError != nil {
		return nil, responseError
	}
	return body, responseError
}

/*
func SearchStartData(data, brand string) []byte {
	method := "POST"
	usr := SearchStart{
		UserID:       532936,
		UserLogin:    "s532936",
		UserPassword: "123456z",
		SearchCode:   data,
		search_cross: "on",
		brand:        brand,
	}
	dts_json_usr, err_json_usr := json.Marshal(usr)
	if err_json_usr != nil {
		fmt.Println(err_json_usr)
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("action", "SearchStart")
	_ = writer.WriteField("data", string(dts_json_usr))
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return body
}
*/
