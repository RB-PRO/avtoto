package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	avtoto "github.com/rb-pro/avtoto"
)

func main() {
	// Загрузка данных из файлов. UserId, UserLogin, UserPassword.
	userIdStr, _ := dataFileExample("UserId.txt")
	userIdInt, _ := strconv.Atoi(userIdStr)
	UserLoginStr, _ := dataFileExample("UserLogin.txt")
	UserPasswordStr, _ := dataFileExample("UserPassword.txt")

	user := avtoto.User{UserId: userIdInt, UserLogin: UserLoginStr, UserPassword: UserPasswordStr} // Объявление пользователя

	mySearchCode := "N007603010406" // Тестовый артикул для поиска

	// ************************** GetBrandsByCode ************************** Поиск бренда по артиклу
	// Создаём структуру запроса бренда по заданному артиклу
	myBrand := avtoto.GetBrandsByCodeRequest{SearchCode: mySearchCode}
	// Получаем с сервера список брендов
	dataGetBrandsByCode, errorSearch := user.GetBrandsByCode(myBrand)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println("> Для артикла", mySearchCode, "найдено", len(dataGetBrandsByCode.Brands), "бренда(ов).",
		"\nПервый найденный бренд имеет производителя", dataGetBrandsByCode.Brands[0].Manuf, "и имя", dataGetBrandsByCode.Brands[0].Name)

	// ************************** SearchStart ************************** Запуск поиска и получение кода ProcessSearchID
	// Объявление запроса метода SearchStart
	searchStartReq := avtoto.SearchStartRequest{SearchCode: mySearchCode, SearchCross: "on", Brand: dataGetBrandsByCode.Brands[0].Manuf}
	// Вызов метода SearchStartRequest с запросом
	datasSearchStartRequest, errorSearch := user.SearchStartRequest(searchStartReq)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println("> Полученный ProcessSearchID", datasSearchStartRequest.ProcessSearchID)

	// ************************** SearchGetParts2 ************************** По коду ProcessSearchID получение найденных данных
	// Преобразовать Ответ метода SearchStart в запрос для метода SearchGetParts2
	SearchGetParts2Req, errorSearch := datasSearchStartRequest.SearchResInReq()
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}

	time.Sleep(1 * time.Second) // Задержка, чтобы сервис нашёл данные на сервере
	// Ответ сервера на запрос
	var SearchGetParts2Res avtoto.SearchGetParts2Response
	for { // В цикле опрашиваем по методу SearchGetParts2 с переданным параметром ProcessSearchID
		// Вызов метода SearchGetParts2
		SearchGetParts2Res, errorSearch = SearchGetParts2Req.SearchGetParts2()
		if errorSearch != nil {
			log.Fatal(errorSearch)
		}

		if SearchGetParts2Res.Error() != "Запрос в обработке" {
			break
		} else {
			fmt.Println("Запрос в обработке. Ждём 1 секунду и заново опрашиваешь по методу SearchGetParts2")
		}
		time.Sleep(1 * time.Second) // Задержка, чтобы сервис нашёл данные на сервере
	}

	fmt.Println("> Всего найдено", len(SearchGetParts2Res.Parts), "деталей, например первая найденная деталь:",
		"\nКод детали", SearchGetParts2Res.Parts[0].Code,
		"\nПроизводитель", SearchGetParts2Res.Parts[0].Manuf,
		"\nНазвание", SearchGetParts2Res.Parts[0].Name,
		"\nЦена", SearchGetParts2Res.Parts[0].Price,
		"\nСклад", SearchGetParts2Res.Parts[0].Storage,
		"\nСрок доставки", SearchGetParts2Res.Parts[0].Delivery,
		"\nМаксимальное количество для заказа", SearchGetParts2Res.Parts[0].MaxCount,
		"\nКратность заказа", SearchGetParts2Res.Parts[0].BaseCount,
		"\nДата обновления склада", SearchGetParts2Res.Parts[0].StorageDate,
		"\nПроцент успешных закупок из общего числа заказов", SearchGetParts2Res.Parts[0].DeliveryPercent,
		"\nПроцент удержания при возврате товара", SearchGetParts2Res.Parts[0].BackPercent,
		"\nНомер запчасти в списке результата поиска", SearchGetParts2Res.Parts[0].AvtotoData.PartId,
		"\nСтатус:", SearchGetParts2Res.Status(),
		"\nSearchID", SearchGetParts2Res.Info.SearchID.Value())

	// ************************** AddToBasket ************************** Добавление товара в корзину
	basketItems := make([]avtoto.AddToBasketRequest, 1)
	basketItem, errorBasketItem := SearchGetParts2Res.SearchResInBasketReq(0)
	if errorBasketItem != nil {
		fmt.Println(errorBasketItem)
	}
	basketItems[0] = basketItem
	basketItems[0].RemoteID = 1
	basketItems[0].Count = 20

	AddToBasketRes, errorRes := user.AddToBasket(basketItems)
	if errorRes != nil {
		fmt.Println(errorRes)
	}

	basketRemoteID := AddToBasketRes.DoneInnerID[0].RemoteID
	basketInnerID := AddToBasketRes.DoneInnerID[0].InnerID
	fmt.Println("> Метод AddToBasket добавил в корзину товар с RemoteID", basketRemoteID, "и InnerID", basketInnerID)

	// ************************** UpdateCountInBasket ************************** Обновление количества товара в корзине
	basketItemsUpdates := make([]avtoto.UpdateCountInBasketRequest, 1)
	basketItemsUpdate, errorBasketItemUpdate := AddToBasketRes.BasketResInUpdateReq(0)
	if errorBasketItemUpdate != nil {
		fmt.Println(errorBasketItemUpdate)
	}
	basketItemsUpdates[0] = basketItemsUpdate
	basketItemsUpdates[0].NewCount = 300

	UpdateCountinBasketRes, errorBasketUpdate := user.UpdateCountInBasket(basketItemsUpdates)
	if errorBasketUpdate != nil {
		fmt.Println(errorBasketUpdate)
	}
	fmt.Println("> Метод UpdateCountinBasketRes выполнился верно для объектов в корзине с RemoteID", UpdateCountinBasketRes.Done)

	// ************************** CheckAvailabilityInBasket ************************** Получить информацию по товару из корзины
	basketChecks := make([]avtoto.CheckAvailabilityInBasketRequest, 1)
	basketCheck, errorbasketChecks := AddToBasketRes.BasketResInCheckReq(0)
	if errorbasketChecks != nil {
		fmt.Println(errorbasketChecks)
	}
	basketChecks[0] = basketCheck

	CheckAvailabilityInBasketRes, errorCheckInBasket := user.CheckAvailabilityInBasket(basketChecks)
	if errorCheckInBasket != nil {
		fmt.Println(errorCheckInBasket)
	}
	availability, errorAvailability := CheckAvailabilityInBasketRes.Availability(0)
	if errorAvailability != nil {
		fmt.Println(errorAvailability)
	}
	fmt.Println("> Метод CheckAvailabilityInBasket.", availability+".", "Максимальное количество товара", CheckAvailabilityInBasketRes.PartsInfo[0].MaxCount)

	/*
		// ************************** AddToOrdersFromBasket ************************** Добавить запчасть из корзины в заказы
		orderBaskets := make([]avtoto.AddToOrdersFromBasketRequest, 1)
		orderBasket, errorbasketChecks := AddToBasketRes.BasketResInOrdersReq(0)
		if errorbasketChecks != nil {
			fmt.Println(errorbasketChecks)
		}
		orderBaskets[0] = orderBasket

		AddToOrdersFromBasketRes, errorOrders := user.AddToOrdersFromBasket(orderBaskets)
		if errorOrders != nil {
			fmt.Println(errorOrders)
		}
		fmt.Println("> Метод AddToOrdersFromBasket.", AddToOrdersFromBasketRes)
	*/

	/*
		// ************************** GetOrdersStatus ************************** Статус заказа
		orderStatusGets := make([]avtoto.GetOrdersStatusRequest, 1)
		orderStatusGet, errorbasketChecks := AddToBasketRes.BasketResInOrdersStatusReq(0)
		if errorbasketChecks != nil {
			fmt.Println(errorbasketChecks)
		}
		orderStatusGets[0] = orderStatusGet

		GetOrdersStatusRes, errorOrdersStatus := user.GetOrdersStatus(orderStatusGets)
		if errorOrdersStatus != nil {
			fmt.Println(errorOrdersStatus)
		}
		orderStatus, orderStatusError := GetOrdersStatusRes.Status(0)
		if orderStatusError != nil {
			fmt.Println(orderStatusError)
		}
		fmt.Println("> Метод GetOrdersStatus.", orderStatus+".", GetOrdersStatusRes.OrdersInfo[0].Info.Progress_text+".", "Всего количество заказов", GetOrdersStatusRes.OrdersInfo[0].Info.Count)
	*/

	// ************************** DeleteFromBasket ************************** Удалить товар из корзины
	basketItemsDeletes := make([]avtoto.DeleteFromBasketRequest, 1)
	basketItemsDelete, errorBasketItemDelete := AddToBasketRes.BasketResInDeleteReq(0)
	if errorBasketItemDelete != nil {
		fmt.Println(errorBasketItemDelete)
	}
	basketItemsDeletes[0] = basketItemsDelete

	DeleteFromBasketRes, errorBusketDelete := user.DeleteFromBasket(basketItemsDeletes)
	if errorBasketItemDelete != nil {
		fmt.Println(errorBusketDelete)
	}
	fmt.Println("> Метод DeleteFromBasketRes выполнился со статусом", DeleteFromBasketRes.Done)

	// ************************** GetStatSearch ************************** статистика проценок по всем объединенным регистрациям.
	statSearch, statSearchError := user.GetStatSearch()
	if statSearchError != nil {
		fmt.Println(statSearchError)
	}

	fmt.Println("> Метод StatSearch вернул информацию о запросах брендов по коду от", statSearch.BrandsStatInfo.StatDateStart.String(), "до", statSearch.BrandsStatInfo.StatDateEndStamp.String(),
		"было", statSearch.BrandsStatInfo.SearchCount, "запроса(ов)")

	// ************************** GetShippingList ************************** получение списка отгрузок.
	GetShippingListReq := avtoto.GetShippingListRequest{}
	ShippingList, ShippingListError := user.GetShippingList(GetShippingListReq)
	if ShippingListError != nil {
		fmt.Println(ShippingListError)
	}
	fmt.Println("> Метод GetShippingList -", ShippingList.Error())
}

// Получение значение из файла
func dataFileExample(filename string) (string, error) {
	// Открыть файл
	fileToken, errorToken := os.Open("example/" + filename)
	if errorToken != nil {
		return "", errorToken
	}

	// Прочитать значение файла
	data := make([]byte, 64)
	n, err := fileToken.Read(data)
	if err == io.EOF { // если конец файла
		return "", errorToken
	}
	fileToken.Close() // Закрытие файла

	return string(data[:n]), nil
}
