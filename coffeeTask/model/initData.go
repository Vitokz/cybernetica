package model
import(
  "flag"
)

type InitDate struct{
  Balance int//:=flag.Int("cash",390,"Количество денег на счету")
  Water int//:=flag.Int("water",540,"Миллилитров воды макс 5л.")
  Milk int//:=flag.Int("milk",400,"Миллилитров молока макс 1л")
  CoffeCorn int //:=flag.Int("corn",120,"Макс 900г")
  CupCount int//:=flag.Int("cup",9,"Количество стаканчиков")
}

func (iD *InitDate) Init() {

    flag.IntVar(&iD.Balance,"cash",390,"Количество денег на счету")
    flag.IntVar(&iD.Water,"water",540,"Миллилитров воды макс 5л.")
    flag.IntVar(&iD.Milk,"milk",400,"Миллилитров молока макс 1л")
    flag.IntVar(&iD.CoffeCorn,"corn",120,"Макс 900г")
    flag.IntVar(&iD.CupCount,"cup",9,"Количеств о стаканчиков")
  flag.Parse()
}