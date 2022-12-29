package avtotoGo

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// Запрос с параметром action и данными json в формате []byte
func HttpPost(bytesRepresentation []byte, action string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("action", action)
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
