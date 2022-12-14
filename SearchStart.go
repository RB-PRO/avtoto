package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func encode_SearchStart(data []byte) ProcessSearchId {
	var result ProcessSearchId
	jsonErr := json.Unmarshal(data, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}
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
