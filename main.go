package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

const url string = "https://www.avtoto.ru/?soap_server=json_mode"

type BrandsByCode struct {
	UserID       int    `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserPassword string `json:"user_password"`
	SearchCode   string `json:"search_code"`
}
type SearchStart struct {
	UserID       int    `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserPassword string `json:"user_password"`
	SearchCode   string `json:"search_code"`
	search_cross string `json:"search_cross"`
	brand        string `json:"brand"`
}
type ProcessSearchId struct {
	ProcessSearchId string `json:"ProcessSearchId"`
}

// ***

type BrandsByCodeReq struct {
	Brands []struct {
		Manuf string `json:"Manuf"`
		Name  string `json:"Name"`
	} `json:"Brands"`
	Info struct {
		Errors       []interface{} `json:"Errors"`
		Logs         []interface{} `json:"Logs"`
		SearchParams struct {
			Number string `json:"number"`
		} `json:"SearchParams"`
		DocVersion string  `json:"DocVersion"`
		Time       float64 `json:"Time"`
	} `json:"Info"`
}
type Seach struct {
	Parts []struct {
		Code            string `json:"Code"`
		Manuf           string `json:"Manuf"`
		Name            string `json:"Name"`
		Price           int    `json:"Price"`
		Storage         string `json:"Storage"`
		Delivery        string `json:"Delivery"`
		MaxCount        string `json:"MaxCount"`
		BaseCount       string `json:"BaseCount"`
		StorageDate     string `json:"StorageDate"`
		DeliveryPercent int    `json:"DeliveryPercent"`
		BackPercent     int    `json:"BackPercent"`
		AvtotoData      struct {
			PartID int `json:"PartId"`
		} `json:"AvtotoData"`
	} `json:"Parts"`
	Info struct {
		SearchID     string        `json:"SearchId"`
		Errors       []interface{} `json:"Errors"`
		Logs         string        `json:"Logs",omitempty`
		SearchParams struct {
			Number  string `json:"number"`
			Analogs int    `json:"analogs"`
		} `json:"SearchParams"`
		DocVersion   string  `json:"DocVersion"`
		Time         float64 `json:"Time"`
		Cache        int     `json:"Cache"`
		SearchStatus int     `json:"SearchStatus"`
	} `json:"Info"`
}

// ***

type requests struct {
	search_cross string
	brand        string
	code         ProcessSearchId
}

func main() {

	//fmt.Println(os.Args)

	var filenameXlsx string
	if len(os.Args) == 1 {
		filenameXlsx = "article.xlsx"
	}
	dataRequest := requestsExcel(filenameXlsx)

	for ind, val := range dataRequest {
		byteValue_start := SearchStartData(val.search_cross, val.brand)
		dataRequest[ind].code = encode_SearchStart(byteValue_start)
	}
	fmt.Println(dataRequest)

	//time.Sleep(8 * time.Second)
	for len(dataRequest) != 0 {
		for ind, val := range dataRequest {
			byteValue_start := SearchGetParts2Data(val.code)
			seach := encode_SearchGetParts2(byteValue_start)
			fmt.Println(seach.Info.Logs)
			if seach.Info.Logs == "" {
				fmt.Println(seach.Parts[0])
				dataRequest = removeByRequests(dataRequest, ind)
			}
		}
		fmt.Println("-> WAIT 3 sec")
		time.Sleep(3 * time.Second)
	}
	/*
		byteValue_start = SearchGetParts2Data(byteValue_start_json)
		//fmt.Println(string(byteValue_start))
		seach := encode_SearchGetParts2(byteValue_start)
		fmt.Println(seach.Parts[0])
	*/

}

func removeByRequests(array []requests, index int) []requests {
	return append(array[:index], array[index+1:]...)
}

func requestsExcel(filename string) []requests {
	var req []requests

	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(f.WorkBook.Sheets.Sheet[0].Name)
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range rows {
		req = append(req, requests{search_cross: row[0], brand: row[1]})
	}

	return req
}
