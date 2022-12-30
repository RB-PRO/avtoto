// Пакет бизнес-логики, согласно ТЗ.

package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/rb-pro/avtoto/pkg/avtotoGo"
)

func Run() {
	// Загрузка данных из файлов. UserId, UserLogin, UserPassword.
	userIdStr, _ := dataFile("UserId.txt")
	userIdInt, _ := strconv.Atoi(userIdStr)
	UserLoginStr, _ := dataFile("UserLogin.txt")
	UserPasswordStr, _ := dataFile("UserPassword.txt")

	user := avtotoGo.User{UserId: userIdInt, UserLogin: UserLoginStr, UserPassword: UserPasswordStr}

	searchStartReq := avtotoGo.SearchStartRequestStruct{SearchCode: "N007603010406", SearchCross: "on", Brand: "MERCEDES-BENZ"}
	jsons, errorSearch := user.SearchStartRequest(searchStartReq)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println(jsons)
}

// Получение значение из файла
func dataFile(filename string) (string, error) {
	fileToken, errorToken := os.Open(filename)
	if errorToken != nil {
		return "", errorToken
	}

	data := make([]byte, 64)
	n, err := fileToken.Read(data)
	if err == io.EOF { // если конец файла
		return "", errorToken
	}
	fileToken.Close()

	return string(data[:n]), nil
}

/*
case 'cannot create client': return 'Не получилось соединиться с сервером';

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
