package kata

import(
  "strings"
  "strconv"
)

func MaxRot(n int64) int64 {
    arr:=strings.Split(strconv.Itoa(int(n)),"")
  
    for i:=0;i<len(arr);i++{
      elem:=arr[i]
      
      for a:=i;a<len(arr)-1;a++{
        arr[a]=arr[a+1]
      }
      
      arr[len(arr)-1]=elem
      str,_:=strconv.Atoi(strings.Join(arr,""))
      if(str>int(n)){
        n=int64(str)
      }
      
    }
  return n
}
