package app

// Файл для работы с методами:
// - GetBrandsByCode
// - SearchStart
// - SearchGetParts2

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

	user := avtotoGo.User{UserId: userIdInt, UserLogin: UserLoginStr, UserPassword: UserPasswordStr} // Объявление пользователя

	mySearchCode := "N007603010406" // Тестовый артикул для поиска

	// ************************** GetBrandsByCode ************************** Поиск бренда по артиклу
	// Создаём структуру запроса бренда по заданному артиклу
	myBrand := avtotoGo.GetBrandsByCodeRequest{SearchCode: mySearchCode}
	// Получаем с сервера список брендов
	dataGetBrandsByCode, errorSearch := user.GetBrandsByCode(myBrand)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println("Для артикла", mySearchCode, "найдено", len(dataGetBrandsByCode.Brands), "бренда(ов).",
		"\nПервый найденный бренд имеет производителя", dataGetBrandsByCode.Brands[0].Manuf, "и имя", dataGetBrandsByCode.Brands[0].Name)

	// ************************** SearchStart ************************** Запуск поиска и получение кода ProcessSearchID
	// Объявление запроса метода SearchStart
	searchStartReq := avtotoGo.SearchStartRequest{SearchCode: mySearchCode, SearchCross: "on", Brand: dataGetBrandsByCode.Brands[0].Manuf}
	// Вызов метода SearchStartRequest с запросом
	datasSearchStartRequest, errorSearch := user.SearchStartRequest(searchStartReq)
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println("Найденный ProcessSearchID", datasSearchStartRequest.ProcessSearchID)

	time.Sleep(8 * time.Second) // Задержка, чтобы сервис нашёл данные на сервере

	// ************************** SearchGetParts2 ************************** По коду ProcessSearchID получение найденных данных
	// Преобразовать Ответ метода SearchStart в запрос для метода SearchGetParts2
	SearchGetParts2Req, errorSearch := datasSearchStartRequest.SearchResInReq()
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	// Вызов метода SearchGetParts2
	SearchGetParts2Res, errorSearch := SearchGetParts2Req.SearchGetParts2()
	if errorSearch != nil {
		log.Fatal(errorSearch)
	}
	fmt.Println("Всего найдено", len(SearchGetParts2Res.Parts), "деталей, например первая найденная деталь:",
		"\nКод детали", SearchGetParts2Res.Parts[0].Code,
		"\nПроизводитель", SearchGetParts2Res.Parts[0].Manuf,
		"\nНазвание", SearchGetParts2Res.Parts[0].Name,
		"\nЦена", SearchGetParts2Res.Parts[0].Price,
		"\nСклад", SearchGetParts2Res.Parts[0].Storage,
		"\nСрок доставки", SearchGetParts2Res.Parts[0].Delivery,
		"\nМаксимальное количество для заказа", SearchGetParts2Res.Parts[0].MaxCount,
		"\nКратность заказа", SearchGetParts2Res.Parts[0].BaseCount)
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
