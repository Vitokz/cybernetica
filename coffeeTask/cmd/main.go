package main

import (
   "main/internal/machine"
   "fmt"
   "main/model"
   //tea "github.com/charmbracelet/bubbletea"
   //"os"
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