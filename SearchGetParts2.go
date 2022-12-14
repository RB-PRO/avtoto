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

func encode_SearchGetParts2(data []byte) Seach {
	var result Seach
	fmt.Println(string(data))
	jsonErr := json.Unmarshal(data, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}

func SearchGetParts2Data(data ProcessSearchId) []byte {
	method := "POST"
	dts_json_usr, err_json_usr := json.Marshal(data)
	if err_json_usr != nil {
		fmt.Println(err_json_usr)
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("action", "SearchGetParts2")
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
