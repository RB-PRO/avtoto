package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
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

type headers struct {
	name        string
	description string
}

func setHeadDatConsidently(f *excelize.File, headData []headers, datas int, color string) {
	var err error
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#" + color}, Pattern: 1},
	})
	if err != nil {
		fmt.Println(err)
	}

	_ = f.SetCellStyle("Sheet1", CNTN(1, 1+datas*10), CNTN(1, 9+datas*10), style)
	for ind, val := range headData {
		f.SetCellValue("Sheet1", CNTN(1, 1+ind+datas*10), val.name)
		err = f.AddComment("Sheet1", CNTN(1, 1+ind+datas*10), `{"author":"RB_PRO: ","text":"`+val.description+`"}`)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func sethead(f *excelize.File) {

	headData := []headers{
		headers{name: "Code", description: "Код детали"},
		headers{name: "Manuf", description: "Производитель"},
		headers{name: "Name", description: "Название"},
		headers{name: "Price", description: "Цена"},
		headers{name: "Storage", description: "Склад"},
		headers{name: "Delivery", description: "Срок доставки"},
		headers{name: "MaxCount", description: "Максимальное количество для заказа, остаток по складу. Значение -1 - означает много или неизвестно."},
		headers{name: "DeliveryPercent", description: "Процент успешных закупок из общего числа заказов"},
		headers{name: "0.9_Price", description: "90% от цены"},
	}
	setHeadDatConsidently(f, headData, 0, fmt.Sprintf("%X%X%X", 150, 200, 150))
	setHeadDatConsidently(f, headData, 1, fmt.Sprintf("%X%X%X", 200, 150, 150))
	setHeadDatConsidently(f, headData, 2, fmt.Sprintf("%X%X%X", 150, 150, 200))

}

func setheadOfDatas(f *excelize.File, datas int) {
}

func writeDate(f *excelize.File, data Seach) {
	Red := rand.Intn(255)
	Green := rand.Intn(255)
	blue := rand.Intn(255)
	h := fmt.Sprintf("%X%X%X", Red, Green, blue)
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#" + h}, Pattern: 1},
	})
	if err != nil {
		fmt.Println(err)
	}

	if len(data.Parts) >= 1 {
		_ = f.SetCellStyle("Sheet1", CNTN(cout, 1), CNTN(cout, 1), style)
		setDataOneSeach(f, data, cout, 0, 0)
	}
	if len(data.Parts) >= 2 {
		setDataOneSeach(f, data, cout, 1, 1)
	}
	if len(data.Parts) >= 3 {
		setDataOneSeach(f, data, cout, 2, 2)
	}
	cout++
}

func setDataOneSeach(f *excelize.File, data Seach, row, datas, index int) {
	f.SetCellValue("Sheet1", CNTN(row, 1+datas*10), data.Parts[index].Code)
	f.SetCellValue("Sheet1", CNTN(row, 2+datas*10), data.Parts[index].Manuf)
	f.SetCellValue("Sheet1", CNTN(row, 3+datas*10), data.Parts[index].Name)
	f.SetCellValue("Sheet1", CNTN(row, 4+datas*10), data.Parts[index].Price)
	f.SetCellValue("Sheet1", CNTN(row, 5+datas*10), data.Parts[index].Storage)
	f.SetCellValue("Sheet1", CNTN(row, 6+datas*10), data.Parts[index].Delivery)
	f.SetCellValue("Sheet1", CNTN(row, 7+datas*10), data.Parts[index].MaxCount)
	f.SetCellValue("Sheet1", CNTN(row, 8+datas*10), data.Parts[index].DeliveryPercent)
	err := f.SetCellFormula("Sheet1", CNTN(row, 9+datas*10), "=PRODUCT("+CNTN(row, 4+datas*10)+",0.9)")
	if err != nil {
		fmt.Println(err)
	}
}
func CNTN(row, col int) string {
	a, _ := excelize.ColumnNumberToName(col)
	return a + strconv.Itoa(row)
}
func main() {

	//fmt.Println(os.Args)

	var filenameXlsx string
	if len(os.Args) == 1 {
		filenameXlsx = "article2.xlsx"
	}
	if len(os.Args) == 2 {
		fmt.Println(os.Args)
	}

	fmt.Println("Ищу файл", filenameXlsx, "в папке с программой.")
	dataRequest := requestsExcel(filenameXlsx)

	fOut := makeFileXLSX("fileOut.xlsx")
	sethead(fOut)

	for ind, val := range dataRequest {
		byteValue_start := SearchStartData(val.search_cross, val.brand)
		dataRequest[ind].code = encode_SearchStart(byteValue_start)
	}

	//time.Sleep(8 * time.Second)
	for len(dataRequest) != 0 {
		fmt.Println("-> Пауза 3 сек.")
		time.Sleep(3 * time.Second)
		for ind := 0; ind < len(dataRequest); ind++ {
			//for ind, val := range dataRequest {
			byteValue_start := SearchGetParts2Data(dataRequest[ind].code)
			seach := encode_SearchGetParts2(byteValue_start)
			/*
				if _, ok := seach.Info["Errors"]; ok {
					fmt.Println(seach.Info["Errors"])
				}
			*/
			//fmt.Println(seach.Info)
			fmt.Println(seach.Parts)
			if seach.Info["Logs"] != "wait" {
				//fmt.Println(len(seach.Parts))
				//fmt.Println(seach.Parts[0].Code)
				seach = sortSeach(seach)
				seach = filterSeach(seach, dataRequest[ind])
				seach = sortThreeSeach(seach)
				if len(seach.Parts) != 0 {
					writeDate(fOut, seach)
					//fmt.Println(seach.Parts)
					dataRequest = removeByRequests(dataRequest, ind)
					fmt.Println("Осталось", len(dataRequest))
				} else {
					dataRequest = removeByRequests(dataRequest, ind)

				}
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

/*
1. Кол-во дней поставки не более 7
2. Кол-во шт. больше 1 (т.е от 2)
3. Цена в файле результаты с точкой сейчас (нужно что бы с запятой было, удобнее для дальнейшей работы и редактирования)
4. Если цена справа кратно отличается от первой или первых двух ( на 30% и более, ставим ее (большую цену))
Пример: цена 100 130 и 135 (первый столбец будет 130, 100 не берем)
*/

func sortSeach(seach Seach) Seach {
	sort.Slice(seach.Parts, func(i, j int) bool {
		return seach.Parts[i].Price < seach.Parts[j].Price
	})
	return seach
}

func sortThreeSeach(seach Seach) Seach {
	if len(seach.Parts) <= 3 {
		return seach
	} else {
		if seach.Parts[0].Price/7.0*10.0 < seach.Parts[1].Price {
			// Убрать 0
			seach = removeBySeach(seach, 0)
		} else {
			// Утрать всё после 3

		}
		for i := 3; i < len(seach.Parts); i++ {
			seach = removeBySeach(seach, i)
			i--
		}
	}
	return seach
}

// Бизнес-логика
func filterSeach(seach Seach, dataRequest requests) Seach {
	// 1 2
	for i := 0; i < len(seach.Parts); i++ {
		Code := seach.Parts[i].Code
		Manuf := seach.Parts[i].Manuf

		Delivery, _ := strconv.Atoi(seach.Parts[i].Delivery)
		MaxCount, _ := strconv.Atoi(seach.Parts[i].MaxCount)
		// 1. Кол-во дней поставки не более 7
		// 2. Кол-во шт. больше 1 (т.е от 2)
		// Проверка на то что артикул и брэнд как в файле
		//fmt.Println("search_cross", dataRequest.search_cross, "brand", dataRequest.brand, "code", dataRequest.code.ProcessSearchId)

		if (Delivery > 7 || MaxCount <= 1) || /*Кол-во дней поставки не более 7, Кол-во шт. больше 1 (т.е от 2) */
			(Code != dataRequest.search_cross || Manuf != dataRequest.brand) { /*Проверка на то что артикул и брэнд как в файле*/
			seach = removeBySeach(seach, i)
			i--
		}

	}
	return seach
}

func removeByRequests(array []requests, index int) []requests {
	return append(array[:index], array[index+1:]...)
}
func removeBySeach(array Seach, index int) Seach {
	theParts := array.Parts
	theParts = append(theParts[:index], theParts[index+1:]...)
	array.Parts = theParts
	return array
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
