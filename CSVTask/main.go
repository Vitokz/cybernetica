package main

import (
	"Vitokz/CSVTask/internal"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	csvFile:=flag.String("file","csvdata/problems.csv","файл с тестами") //Флаг задающий имя файла который нужно спрасить и провести тест
    seconds:=flag.Int("timer",30,"Время на ответ")
	flag.Parse()
	//Открытие файла CSV
	file, err := os.Open(*csvFile)
	//Проверка открылся ли фаил
	if err != nil {
		panic(err)
	}
	//defer выполняет инструкцию после окончания функции
	//Close закрывает фаил
	defer file.Close()

	reader := csv.NewReader(file)        //NewReader  читает фаил и создает класс который читает данные из file
	lencsv, _ := reader.Read()         //Достаем первую строку для определения длинны строки
	reader.FieldsPerRecord = len(lencsv) //берем элементы до числа len(lencsv)

	var arrayAns = map[string]int{
		"right":     0,
		"incorrect": 0,
	}
	_ = arrayAns
	for {
		record, e := reader.Read() //Достаю строки поочередно
		if e != nil {
			fmt.Println(e)
			break
		}
        
		fmt.Println(record[0] + "=")
		fmt.Printf("У вас %d с. Введите ответ:", *seconds)

		answer := internal.Input(*seconds)

		check,ok:=internal.Check(strings.TrimSpace(answer),record[1])
		fmt.Println(check + "\n")
		if ok {
			arrayAns["right"]++
		}else{
			arrayAns["incorrect"]++
		}
	}
	fmt.Printf("Правильных ответов:%d Неправильных ответов:%d",arrayAns["right"],arrayAns["incorrect"])
}
