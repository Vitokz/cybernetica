package kata

import(
    "math"
)

func NewAvg(arr []float64, navg float64) int64 {
   sum:=float64(len(arr)+1)*navg
  
   sumarr:=float64(0)
  for _,v:=range arr {
    sumarr+=v
  }
   result:=sum-sumarr
  
  if result<=0{
    return -1
  }
  
  return int64(math.Ceil(result))
}
