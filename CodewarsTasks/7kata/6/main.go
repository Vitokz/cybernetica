package kata

import(
"strings"
)

func Capitalize(st string, arr []int) string {
  strarr:=strings.Split(st,"")
  lenstr:=len(strarr)
  for _,v :=range arr {
    if v<=lenstr-1 {
      strarr[v]=strings.ToUpper(strarr[v])
    }
  }
  return strings.Join(strarr,"")
}
