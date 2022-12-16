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
	var err error
	f.SetCellValue("Sheet1", "A1", "Code")
	err = f.AddComment("Sheet1", "A1", `{"author":"RB_PRO: ","text":" Код детали"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "B1", "Manuf")
	err = f.AddComment("Sheet1", "B1", `{"author":"RB_PRO: ","text":" Производитель"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "C1", "Name")
	err = f.AddComment("Sheet1", "C1", `{"author":"RB_PRO: ","text":" Название"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "D1", "Price")
	err = f.AddComment("Sheet1", "D1", `{"author":"RB_PRO: ","text":" Цена"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "E1", "Storage")
	err = f.AddComment("Sheet1", "E1", `{"author":"RB_PRO: ","text":" Склад"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "F1", "Delivery")
	err = f.AddComment("Sheet1", "F1", `{"author":"RB_PRO: ","text":" Срок доставки"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "G1", "MaxCount")
	err = f.AddComment("Sheet1", "G1", `{"author":"RB_PRO: ","text":" Максимальное количество для заказа, остаток по складу. Значение -1 - означает много или неизвестно."}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "H1", "BaseCount")
	err = f.AddComment("Sheet1", "H1", `{"author":"RB_PRO: ","text":" Кратность заказа"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "I1", "StorageDate")
	err = f.AddComment("Sheet1", "I1", `{"author":"RB_PRO: ","text":" Дата обновления склада"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "J1", "DeliveryPercent")
	err = f.AddComment("Sheet1", "J1", `{"author":"RB_PRO: ","text":" Процент успешных закупок из общего числа заказов"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "K1", "BackPercent")
	err = f.AddComment("Sheet1", "K1", `{"author":"RB_PRO: ","text":" Процент удержания при возврате товара (при отсутствии возврата поставщику возвращается значение -1)"}`)
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue("Sheet1", "L1", "0.9_Price")
	err = f.AddComment("Sheet1", "L1", `{"author":"RB_PRO: ","text":" 90% от цены"}`)
	if err != nil {
		fmt.Println(err)
	}
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

		err := f.SetCellFormula("Sheet1", "L"+strconv.Itoa(cout), "=PRODUCT("+"D"+strconv.Itoa(cout)+",0.9)")
		if err != nil {
			fmt.Println(err)
		}
		cout++
	}
}

func main() {

	//fmt.Println(os.Args)

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
	sethead(fOut)

	for ind, val := range dataRequest {
		byteValue_start := SearchStartData(val.search_cross, val.brand)
		dataRequest[ind].code = encode_SearchStart(byteValue_start)
	}
	//fmt.Println(dataRequest)

	//time.Sleep(8 * time.Second)
	for len(dataRequest) != 0 {
		fmt.Println("-> Пауза 3 секунд")
		time.Sleep(3 * time.Second)
		for ind := 0; ind < len(dataRequest); ind++ {
			//for ind, val := range dataRequest {
			byteValue_start := SearchGetParts2Data(dataRequest[ind].code)
			seach := encode_SearchGetParts2(byteValue_start)
			if seach.Info["Logs"] != "wait" {
				//fmt.Println(seach.Parts[0].Code)
				seach = filterSeach(seach)
				writeDate(fOut, seach)
				//fmt.Println("Len", len(dataRequest), "ind", ind)
				dataRequest = removeByRequests(dataRequest, ind)
				fmt.Println("Осталось", len(dataRequest))
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

// Бизнес-логика
func filterSeach(seach Seach) Seach {
	var out Seach
	// 1 2
	for i := 0; i < len(seach.Parts); i++ {
		code, _ := strconv.Atoi(seach.Parts[i].Code)
		MaxCount, _ := strconv.Atoi(seach.Parts[i].MaxCount)
		if code <= 7 && MaxCount > 1 {
			out.Parts = append(out.Parts, seach.Parts[i])
		}
	}
	return out
}

func removeByRequests(array []requests, index int) []requests {
	return append(array[:index], array[index+1:]...)
}
func removeBySeach(array []Seach, index int) []Seach {
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
