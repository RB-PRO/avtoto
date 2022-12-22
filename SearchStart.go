package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type SearchStartRequest struct {
	user_id       int    `json:"user_id"`       // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	user_login    string `json:"user_login"`    // Логин пользователя (тип: строка)
	user_password string `json:"user_password"` // Пароль пользователя (тип: строка)
	search_code   string `json:"search_code"`   // Поисковый запрос, минимум 3 символа (тип: строка)
	search_cross  string `json:"search_cross"`  // Искать в аналогах или нет (тип: строка, 'on' или 'off')
	brand         string `json:"brand"`         /// Искать код с учетом бренда, минимум 2 символа (опционально)(тип: строка)
	// [*] — эти данные можно узнать зайдя на страницу Настройки после авторизации на сайте
	// [**] Список брендов можно получить с помощью метода GetBrandsByCode
}

type searchStartRequestData struct {
	action string             `json:"action"`
	data   SearchStartRequest `json:"data"`
}

func (searchReq SearchStartRequest) Post() {
	RequestData := searchStartRequestData{action: "SearchStart", data: searchReq}

	personJSON, _ := json.Marshal(RequestData)

	response, error := http.Post(URL,
		"data", bytes.NewBuffer(personJSON))

	if error != nil {
		print(error)
	}

	defer response.Body.Close()
	body, errorasd := ioutil.ReadAll(response.Body)
	if errorasd != nil {
		print(errorasd)
	}

	fmt.Println(string(body))
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
