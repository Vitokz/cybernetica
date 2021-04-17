package internal

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func Check(userAns, rightAns string) (string, bool) {
	rand.Seed(time.Now().UnixNano())
	right := []string{"Правильно!", "Лучший", "Повезло Повезло", "Бог математики"}
	incorrect := []string{"Ну почти", "Неправильно", "даже Сосик из 7б решил", "Пифагор осуждает"}
	usAns, err := strconv.Atoi(userAns)
	if err != nil {
		fmt.Println("Буквы вводить запрещено")
		return incorrect[rand.Intn(4)], false
	}
	riAns, _ := strconv.Atoi(rightAns)
	if err != nil {
		log.Fatal(err.Error())
	}

	if usAns == riAns {
		return right[rand.Intn(4)], true
	}
	return incorrect[rand.Intn(4)], false
}
