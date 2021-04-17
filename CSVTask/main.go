package main

import (
	"Vitokz/CSVTask/internal"
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//  Флаги
	/*Флаг для вставки файла*/
	csvFile := flag.String("file", "csvdata/problems.csv", "файл с тестами") //Флаг задающий имя файла который нужно спрасить и провести тест
	/*флаг изменения времени*/
	seconds := flag.Int("timer", 30, "Время на ответ")
	/*Парсинг флагов*/
	flag.Parse()
	//

	//Открытие файла CSV
	file, err := os.Open(*csvFile)
	//Проверка открылся ли фаил
	if err != nil {
		panic(err)
	}
	//defer выполняет инструкцию после окончания функции
	//Close закрывает фаил
	defer file.Close()

	reader := csv.NewReader(file) //NewReader  читает фаил и создает класс который читает данные из file
	reader.FieldsPerRecord = 2    //берем элементы до 2 каждой строки

	//Массив для сбора информации о правильных и неправильных ответах
	var arrayAns = map[string]int{
		"right":     0,
		"incorrect": 0,
	}

	/* Канал */
	chanel := make(chan string)
	/* Таймер */ timer := time.NewTimer(time.Duration(*seconds) * time.Second)
loop: // "Имя" функции
	for {
		record, e := reader.Read() //Достаю строки поочередно
		if e != nil {              //Проверка досталась ли строка
			break
		}

		fmt.Print(record[0] + "=")
		fmt.Print("Введите ответ:")

		go func() { //Горутина для ввода ответов
			reader := bufio.NewReader(os.Stdin)
			inp, _ := reader.ReadString('\n')
			chanel <- inp
		}()

		select {
		case <-timer.C: //Кейс проверяющий не остановился ли таймер
			break loop
		case answer := <-chanel: //Засчитывание ответов
			check, ok := internal.Check(strings.TrimSpace(answer), record[1]) //Функция проверки ответа
			fmt.Println(check + "\n")
			if ok {
				arrayAns["right"]++
			} else {
				arrayAns["incorrect"]++
			}
		}
	}

	fmt.Println("Правильных ответов:" + strconv.Itoa(arrayAns["right"]) + " Неправильных ответов:" + strconv.Itoa(arrayAns["incorrect"]))
}
