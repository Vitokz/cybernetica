package kata

import(
  "strings"
  "strconv"
)
func HighAndLow(in string) string {
  var max int=-10000
  var min int=10000
  numbers := strings.Split(in," ")
  for i :=0;i < len(numbers);i++{
    nm,err:=strconv.Atoi(numbers[i])
    if(err != nil){
      return "ОШИБОЧКА"
    }
    
    if (nm<min){
      min=nm
    }
    if(nm>max){
      max=nm
    }
  }
  return strconv.Itoa(max) + " " + strconv.Itoa(min)
}
