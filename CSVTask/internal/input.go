package internal

import (
	"bufio"
	"time"
	"os"
)


func Input(seconds int) string{
_=seconds
timer := time.NewTimer(time.Second)
<-timer.C
answer := bufio.NewReader(os.Stdin)
time.Sleep(time.Duration(seconds))


inp, _ := answer.ReadString('\n')
return inp
}