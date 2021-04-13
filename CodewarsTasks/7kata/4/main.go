package kata


func Seven(n int64) []int {
    count:=0
    for n>=100 {
      n=n/10-2*(n%10)
      count++
    }
    return []int{int(n),count}
}
