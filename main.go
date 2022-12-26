package main

import (
	"fmt"
	"log"
)

type User struct {
	UserId       int    `json:"user_id"`       // Уникальный идентификатор пользователя (номер клиента) (тип: целое)
	UserLogin    string `json:"user_login"`    // Логин пользователя (тип: строка)
	UserPassword string `json:"user_password"` // Пароль пользователя (тип: строка)
}

const URL string = "https://www.avtoto.ru/?soap_server=json_mode"

func main() {
	//sreq := SearchStartRequestStruct{user_id: 532936, user_login: "s532936", user_password: "123456z", search_code: "N007603010406", search_cross: "off", brand: "MERCEDES-BENZ"}
	//sreq.Post()
	user := User{UserId: 532936, UserLogin: "s532936", UserPassword: "123456z"}
	searchStartReq := SearchStartRequestStruct{SearchCode: "N007603010406", SearchCross: "on", Brand: "MERCEDES-BENZ"}
	jsons, errorSearch := user.SearchStartRequest(searchStartReq)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println(jsons)
}

/*case 'cannot create client': return 'Не получилось соединиться с сервером';
  case 'no result':            return 'Сервер не ответил';
  case 'wrong params':         return 'Неверные параметры соединения';
  case 'wrong parts':          return 'Ошибка данных';
  case 'error code':           return 'Неверный артикул';
  private $progress_list = array(
      '2'=>  'Ожидает оплаты',
      '1'=>  'Ожидает обработки',
      '3'=>  'Заказано',
      '4'=>  'Закуплено',
      '5'=>  'В пути',
      '6'=>  'На складе',
      '7'=>  'Выдано',
      '8'=>  'Нет в наличии'
  );
*/
