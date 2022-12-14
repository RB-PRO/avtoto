package main

import (
	"fmt"
	"os"
	"strconv"
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
	Info map[string]interface{} `json:"Info"`
	/*
		Info struct {
			SearchID     string        `json:"SearchId"`
			Errors       []interface{} `json:"Errors"`
			Logs         string        `json:"Logs,omitempty"`
			SearchParams struct {
				Number  string `json:"number"`
				Analogs int    `json:"analogs"`
			} `json:"SearchParams"`
			DocVersion   string  `json:"DocVersion"`
			Time         float64 `json:"Time"`
			Cache        int     `json:"Cache"`
			SearchStatus int     `json:"SearchStatus"`
		} `json:"Info"`
	*/
}

// ***

type requests struct {
	search_cross string
	brand        string
	code         ProcessSearchId
}

func makeFileXLSX(filename string) *excelize.File {
	f := excelize.NewFile()
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
	return f
}

var cout int = 2

func sethead(f *excelize.File) {
	f.SetCellValue("Sheet1", "A1", "Code")
	f.SetCellValue("Sheet1", "B1", "Manuf")
	f.SetCellValue("Sheet1", "C1", "Name")
	f.SetCellValue("Sheet1", "D1", "Price")
	f.SetCellValue("Sheet1", "E1", "Storage")
	f.SetCellValue("Sheet1", "F1", "Delivery")
	f.SetCellValue("Sheet1", "G1", "MaxCount")
	f.SetCellValue("Sheet1", "H1", "BaseCount")
	f.SetCellValue("Sheet1", "I1", "StorageDate")
	f.SetCellValue("Sheet1", "J1", "DeliveryPercent")
	f.SetCellValue("Sheet1", "K1", "BackPercent")
}

func writeDate(f *excelize.File, data Seach) {
	for _, val := range data.Parts {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(cout), val.Code)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(cout), val.Manuf)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(cout), val.Name)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(cout), val.Price)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(cout), val.Storage)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(cout), val.Delivery)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(cout), val.MaxCount)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(cout), val.BaseCount)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(cout), val.StorageDate)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(cout), val.DeliveryPercent)
		f.SetCellValue("Sheet1", "K"+strconv.Itoa(cout), val.BackPercent)
		cout++
	}
}

func main() {

	fmt.Println(os.Args)

	var filenameXlsx string
	if len(os.Args) == 1 {
		filenameXlsx = "article.xlsx"
	}
	if len(os.Args) == 2 {
		fmt.Println(os.Args)
	}

	fmt.Println("Ищу файл", filenameXlsx, "в папке с программой.")
	dataRequest := requestsExcel(filenameXlsx)

	fOut := makeFileXLSX("fileOut.xlsx")

	for ind, val := range dataRequest {
		byteValue_start := SearchStartData(val.search_cross, val.brand)
		dataRequest[ind].code = encode_SearchStart(byteValue_start)
	}
	//fmt.Println(dataRequest)

	//time.Sleep(8 * time.Second)
	for len(dataRequest) != 0 {
		fmt.Println("-> Пауза 5 секунд")
		time.Sleep(5 * time.Second)
		for ind := 0; ind < len(dataRequest); ind++ {
			//for ind, val := range dataRequest {
			byteValue_start := SearchGetParts2Data(dataRequest[ind].code)
			seach := encode_SearchGetParts2(byteValue_start)
			if seach.Info["Logs"] != "wait" {
				//fmt.Println(seach.Parts[0].Code)
				writeDate(fOut, seach)
				//fmt.Println("Len", len(dataRequest), "ind", ind)
				dataRequest = removeByRequests(dataRequest, ind)
				fmt.Println("Осталось:", len(dataRequest))
			}
		}
	}
	fOut.Save()
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
