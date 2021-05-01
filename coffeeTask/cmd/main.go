package main

import (
   "main/internal/machine"
   "fmt"
   "main/model"
)

func main() {
  fl:=model.InitDate{}
  fl.Init()

  mach, err := machine.New(&fl)
  if err != nil {
    fmt.Println(err)
  }
  mach.Start()
}