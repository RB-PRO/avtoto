package app

// Файл для работы с методами:
// - AddToBasket
// - UpdateCountInBasket
// - DeleteFromBasket
// - CheckAvailabilityInBasket

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rb-pro/avtoto/pkg/avtotoGo"
)

func Basket() {
	// Загрузка данных из файлов. UserId, UserLogin, UserPassword.
	userIdStr, _ := dataFile("UserId.txt")
	userIdInt, _ := strconv.Atoi(userIdStr)
	UserLoginStr, _ := dataFile("UserLogin.txt")
	UserPasswordStr, _ := dataFile("UserPassword.txt")

	user := avtotoGo.User{UserId: userIdInt, UserLogin: UserLoginStr, UserPassword: UserPasswordStr} // Объявление пользователя

	/*
		Код детали N007603010406
		Производитель КИТАЙ
		Название Кольцо уплотнительное Mercedes-Benz N007603010406 (10)
		Цена 79
		Склад Москва
		Срок доставки 17
		Максимальное количество для заказа 37000
		Кратность заказа 10
	*/
	/*{"Code":"N007603010406","Manuf":"КИТАЙ","Name":"Кольцо уплотнительное Mercedes-Benz N007603010406 (10)","Price":79,"Storage":"Москва","Delivery":"17","MaxCount":"37000","BaseCount":"10","StorageDate":"31.12.2022","DeliveryPercent":86,"BackPercent":-1,"AvtotoData":{"PartId":0}}*/
	/*"SearchId":116526597*/
	AddToBasketReq := make([]avtotoGo.AddToBasketRequest, 1)
	AddToBasketReq[0].Code = "N007603010406"
	AddToBasketReq[0].Manuf = "КИТАЙ"
	AddToBasketReq[0].Name = "Кольцо уплотнительное Mercedes-Benz N007603010406 (10)"
	AddToBasketReq[0].Storage = "Москва"
	AddToBasketReq[0].Delivery = "17"
	//AddToBasketReq[0].PartId = "0"
	//AddToBasketReq[0].SearchID = "116526597"

	AddToBasketRes, errorBasket := user.AddToBasket(AddToBasketReq)
	if errorBasket != nil {
		log.Println(errorBasket)
	}
	fmt.Println(AddToBasketRes)
}
