// Пакет бизнес-логики, согласно ТЗ.

package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/rb-pro/avtoto/pkg/avtotoGo"
)

func Run() {
	// Загрузка данных из файлов. UserId, UserLogin, UserPassword.
	userIdStr, _ := dataFile("UserId.txt")
	userIdInt, _ := strconv.Atoi(userIdStr)
	UserLoginStr, _ := dataFile("UserLogin.txt")
	UserPasswordStr, _ := dataFile("UserPassword.txt")

	// Объявление пользователя
	user := avtotoGo.User{UserId: userIdInt, UserLogin: UserLoginStr, UserPassword: UserPasswordStr}

	// ************************** SearchStart **************************

	// Объявление запроса метода SearchStart
	searchStartReq := avtotoGo.SearchStartRequest{SearchCode: "N007603010406", SearchCross: "on", Brand: "MERCEDES-BENZ"}
	// Вызов метода SearchStartRequest с запросом
	jsonsSearchStartRequest, errorSearch := user.SearchStartRequest(searchStartReq)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println(jsonsSearchStartRequest)

	// ************************** SearchGetParts2 **************************

	// Преобразовать Ответ метода SearchStart в запрос для метода SearchGetParts2
	SearchGetParts2Req, errorSearch := jsonsSearchStartRequest.SearchResInReq()
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}

	time.Sleep(8 * time.Second)

	// Вызов метода SearchGetParts2
	SearchGetParts2Res, errorSearch := SearchGetParts2Req.SearchGetParts2()
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println(SearchGetParts2Res)
}

// Получение значение из файла
func dataFile(filename string) (string, error) {
	// Открыть файл
	fileToken, errorToken := os.Open(filename)
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
