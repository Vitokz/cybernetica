package repository
import(
  "main/model"
)
var (
  Espresso model.CoffeeModel = model.CoffeeModel{
    Water: 250,
    CoffeCorn: 16,
    Milk:0,
    Price:60,
  }

  Latte model.CoffeeModel= model.CoffeeModel{
    Water: 300,
    CoffeCorn:20,
    Milk:76,
    Price:110,  
  }

  Cappuchino model.CoffeeModel = model.CoffeeModel{
    Water:200,
    CoffeCorn:16,
    Milk:100,
    Price:140,
  }
)