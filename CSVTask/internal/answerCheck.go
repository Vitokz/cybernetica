package internal

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
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
	riAns, err := strconv.Atoi(rightAns)
	if err != nil {
		errWrapped := errors.Wrap(err, "Неверные данные в csv")
		fmt.Println(errWrapped)
		os.Exit(1)
	}

	if usAns == riAns {
		return right[rand.Intn(4)], true
	}
	return incorrect[rand.Intn(4)], false
}
