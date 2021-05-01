package machine

import(
   "bufio"
   "strconv"
   "fmt"
   "strings"
   "os"
   "errors"
   "main/model"
   "main/repository"
   "main/proto"
 )
type CoffeeMachine struct { // Струтура кофемашины
  Storage model.StorageModel //Модель расходников машины
  Balance int // касса машины
  Stat model.Stat //статистика заказов
}

func New(init *model.InitDate) (CoffeeMachine,error){
  machine:=CoffeeMachine{}

  setMilk(&machine,*init)
  setCorn(&machine,*init)
  setBalance(&machine,*init)
  setWater(&machine,*init)
  setCup(&machine,*init)
  return machine,nil
}

func fill(m *CoffeeMachine){
       fmt.Print ("Введите пароль")//7788
	       reader:=bufio.NewReader(os.Stdin) //Ввод пароля для проверки на админа
	       inp,_:=reader.ReadString('\n')  
	       if(strings.TrimSpace(inp)!= "7788"){
	         err:=errors.New("Вы посторонний человек уходите")
           fmt.Println(err)
           return
	       }
     fmt.Printf("")
      plusCorn(m)
      plusCup(m)
      plusMilk(m)
      plusWater(m)
}

func (m *CoffeeMachine) Start()  {
  fmt.Println("Добро пожаловать, чего изволите?")
  loop:
     for{  //бесконечный цикл
     fmt.Println("buy- команда заказа кофе")
     fmt.Println("stat- команда вывода статистики с момента запуска")
     fmt.Println("fill- команда команда для админа ,добавляющая ресурсы")
     fmt.Println("take- команда вывода денег из кофемашины")
     fmt.Println("Введите команду: ")
       reader:=bufio.NewReader(os.Stdin) //Ввод команды
       inp,_:=reader.ReadString('\n')  //обработчик ввода
     switch str:=strings.TrimSpace(inp);str{ //Switch вызывающий нужныe инструкции при опреденных вводах
       case "buy":   
        m.buy()
       case "fill":
        m.fill()
       case "take":
       err:=m.take()
       _=err
       case "stat":
       m.stat()
       default:
       if(str=="exit"){
        break loop
       }
        fmt.Println("вы ввели неправильную команду")
     }
     }
     fmt.Println("Досвидания")
}

func (m *CoffeeMachine) buy() {
     fmt.Println("На данный момент на выбор есть 3 кофе:")
     fmt.Println("1-espresso Цена- 60р.")
     fmt.Println("2-latte Цена- 110р.")
     fmt.Println("2-cappuchino Цена- 140р.")
     fmt.Println("Выберете напиток: 1-espresso 2-latte 3-cappuchino")

     reader:=bufio.NewReader(os.Stdin) //Ввод номера напитка
     inp,_:=reader.ReadString('\n')

     switch number:=strings.TrimSpace(inp);number{ //Swwitch вызываюший нужные инструкции определенных напитков

       case "1": 
        err:=checkStorage(repository.Espresso,m.Storage) //Проверяет хватает ли ресурсов машины для создания напитка
          if(err != nil){
            fmt.Println(err)
           fmt.Println("Попробуйте выбрать что-нибудь другое")
            m.buy()
        }
        m.Balance+=repository.Espresso.Price //Плюсуем к кассе сумму заказа
        m.Stat.Espresso+=1  //Добаляем заказ в статистику определенных напитков
        minusStorage(repository.Espresso, &m.Storage) // Минусуем ресурсы которые потребовались для выполнения заказа
        fmt.Println("Через +-10 секунд ваш кофе будет готов")
       case "2": // Комментарии выше актуальны и тут
         err:=checkStorage(repository.Latte,m.Storage)
          if(err != nil){
            fmt.Println(err)
            fmt.Println("Попробуйте выбрать что-нибудь другое")
            m.buy()
          }
        m.Balance+=repository.Latte.Price
        m.Stat.Latte+=1
        minusStorage(repository.Latte, &m.Storage)
        fmt.Println("Через +-10 секунд ваш кофе будет готов")

       case "3":// Комментарии выше актуальны и тут
         err:=checkStorage(repository.Cappuchino,m.Storage)
          if(err != nil){
            fmt.Println(err)
            fmt.Println("Попробуйте выбрать что-нибудь другое")
            m.buy()
          }
        m.Balance+=repository.Cappuchino.Price
        m.Stat.Cappuchino+=1
        minusStorage(repository.Cappuchino, &m.Storage)
        fmt.Println("Через +-10 секунд ваш кофе будет готов")

       default: 
       if(number!="back"){ // в случает ввода непонятно чего предупреждаем об этом клиента
         err:=errors.New("Вы ввели значение,которое не может обрабатываться машиной")
         fmt.Println(err)
         fmt.Println("Даем вам еще один шанс")
         m.buy()
       }
       break
     }
    
 }


 func (m *CoffeeMachine) fill() error {
     fmt.Print ("Введите пароль")//7788
     reader:=bufio.NewReader(os.Stdin) //Ввод пароля для проверки на админа
     inp,_:=reader.ReadString('\n')  
     if(strings.TrimSpace(inp)!= "7788"){
       err:=errors.New("Вы посторонний человек уходите")
       return err
     }
      plusCorn(m)
      plusCup(m)
      plusMilk(m)
      plusWater(m)
    return nil  
 }



func (m *CoffeeMachine) take() error {
     fmt.Print ("Введите пароль: ")//7788
     err:=authInput(0)
     //fmt.Println(err)
     if err!=nil{
       return err
     }
     fmt.Printf("Вы вывели %d путинских дублонов\n", m.Balance) //Вывод и обнуление баланса
     m.Balance=0
     return nil
     
}

func (m *CoffeeMachine) stat()  {
    fmt.Println("В машине на данный момент:")
    fmt.Printf("%d мл воды\n", m.Storage.Water)
    fmt.Printf("%v мл молока\n", m.Storage.Milk)
    fmt.Printf("%v г кофейных зерен\n", m.Storage.CoffeCorn)
    fmt.Printf("%d стаканчиков\n", m.Storage.CupCoount)
    fmt.Printf("А в кассе тем временем %d\n", m.Balance)
    fmt.Println(" ")
    fmt.Println(" ")
    fmt.Printf("%d-espresso\n",m.Stat.Espresso)
    fmt.Printf("%d-latte\n",m.Stat.Latte)
    fmt.Printf("%d-cappuchino\n",m.Stat.Cappuchino)
    sumBuy:=m.Stat.Espresso+m.Stat.Latte+m.Stat.Cappuchino
    cash:=m.Stat.Espresso*repository.Espresso.Price+m.Stat.Latte*repository.Latte.Price+m.Stat.Cappuchino*repository.Cappuchino.Price
    fmt.Printf("Всего напитков продано %d на %d путинских дублонов\n", sumBuy,cash)
}


//Для функции Buy
func checkStorage (name model.CoffeeModel,m model.StorageModel) error{ //Проверка на достаточность ресурсов
  switch{
          case m.CoffeCorn<name.CoffeCorn:
           err:=errors.New("К сожалению в данный момент недостаточно кофе для приготовления")
          return err
          case m.Milk < name.Milk:
          err:=errors.New("К сожалению в данный момент недостаточно молока для приготовления")
          return err
          case m.Water < name.Water:
          err:=errors.New("К сожалению в данный момент недостаточно воды для приготовления")
          return err
  }
  return nil
}
//Для функции Buy
func minusStorage (name model.CoffeeModel,m *model.StorageModel){  // Минусуем затрачиваемые ресурсы
     m.CoffeCorn-=name.CoffeCorn
     m.Milk-=name.Milk
     m.CupCoount--
     m.Water-=name.Water
}
//Относятся к new
func setMilk(m *CoffeeMachine, init model.InitDate){
    if (init.Milk >proto.MILK) {
      may:=strconv.Itoa(proto.MILK)
      err:=errors.New("Вы не можете залить молока больше чем:"+may+"ml")
      fmt.Println(err)
      fmt.Println("Попробуйте еще раз")
      init.Milk=input(proto.MILK)
    }
    m.Storage.Milk=init.Milk
}
func setCorn (m *CoffeeMachine, init model.InitDate){
    if (init.CoffeCorn >proto.CORN) {
      may:=strconv.Itoa(proto.CORN)
      err:=errors.New("Вы не можете засыпать больше чем:"+may+"зерен")
      fmt.Println(err)
      fmt.Println("Попробуйте еще раз")
      init.CoffeCorn=input(proto.CORN)
    }
    m.Storage.CoffeCorn=init.CoffeCorn
}
func setWater (m *CoffeeMachine, init model.InitDate){
    if (init.Water >proto.WATER) {
      may:=strconv.Itoa(proto.WATER)
      err:=errors.New("Вы не можете залить воды больше чем:"+may+"ml")
      fmt.Println(err)
      fmt.Println("Попробуйте еще раз")
      init.Water=input(proto.WATER)
    }
    m.Storage.Water=init.Water
}
func setCup (m *CoffeeMachine, init model.InitDate){
    if (init.CupCount >proto.CUP) {
      may:=strconv.Itoa(proto.CUP)
      err:=errors.New("Вы не можете засунуть стаканчиков больше чем:"+may+"шт")
      fmt.Println(err)
      fmt.Println("Попробуйте еще раз")
      init.CupCount=input(proto.CUP)
    }
    m.Storage.CupCoount=init.CupCount
}
func setBalance (m *CoffeeMachine, init model.InitDate){
    if (init.Balance >proto.BALANCE) {
      may:=strconv.Itoa(proto.BALANCE)
      err:=errors.New("Вы не можете пополнить баланс больше чем:"+may+"шт")
      fmt.Println(err)
      fmt.Println("Попробуйте еще раз")
      init.Balance=input(proto.BALANCE)
    }
    m.Balance=init.Balance
}
//
//Относятся к fill
 func plusWater(m *CoffeeMachine){
     fmt.Printf("На данный момент в машине %v мл воды из %v\n",m.Storage.Water,proto.WATER)
     fmt.Println("Сколько воды вы хотите добавить?")
      water:=plusInput(proto.WATER,m.Storage.Water)
      m.Storage.Water+=water 
 }

 func plusCorn(m *CoffeeMachine){
     fmt.Printf("На данный момент в машине %v г зерна из %v\n",m.Storage.CoffeCorn,proto.CORN)
     fmt.Println("Сколько грамм зерна вы хотите добавить?")
      corn:=plusInput(proto.CORN,m.Storage.CoffeCorn)
      m.Storage.CoffeCorn+=corn 
 }
 func plusMilk(m *CoffeeMachine){
     fmt.Printf("На данный момент в машине %v мл молока из %v\n",m.Storage.Milk,proto.MILK)
     fmt.Println("Сколько воды вы хотите добавить?")
      milk:=plusInput(proto.MILK,m.Storage.Milk)
      m.Storage.Milk+=milk 
 }
 func plusCup(m *CoffeeMachine){
     fmt.Printf("На данный момент в машине %v стаканчиков из %v\n",m.Storage.CupCoount,proto.CUP)
     fmt.Println("Сколько воды вы хотите добавить?")
      cup:=plusInput(proto.CUP,m.Storage.CupCoount)
      m.Storage.CupCoount+=cup
 }

 func input(limit int) int{
      fmt.Println("Введите значение")
      reader:=bufio.NewReader(os.Stdin) 
      inp,_:=reader.ReadString('\n')
      result,err:=strconv.Atoi(strings.TrimSpace(inp))
      if err!=nil{
        fmt.Println("Нельзя вводить буквы вместо чисел! Попробуйте еще раз")
        input(limit)
      }
      if(result>limit){
        may:=strconv.Itoa(proto.MILK)
       fmt.Println("Введите число меньше"+may)
       input(limit)
      }
      return result
}
func plusInput(limit int,real int) int{
  fmt.Println("Введите значение")
      reader:=bufio.NewReader(os.Stdin) 
      inp,_:=reader.ReadString('\n')
      result,err:=strconv.Atoi(strings.TrimSpace(inp))
      if err!=nil{
        fmt.Println("Нельзя вводить буквы вместо чисел! Попробуйте еще раз")
        input(limit)
      }
      if(result>limit-real){
        may:=strconv.Itoa(limit-real)
       fmt.Println("Введите число меньше либо равное: "+may)
       input(limit)
      }
      return result
}
func authInput(count int) error{
  if(count>=3){
    err:=errors.New("Уходи ты не админ")
    return err
  }
     reader:=bufio.NewReader(os.Stdin) //Ввод пароля для проверки на админа
     inp,_:=reader.ReadString('\n')  
     if(strings.TrimSpace(inp)!= "7788"){
       err:=errors.New("Неправильный пароль")
       fmt.Print(err)
       fmt.Printf(", попробуйте снова,осталось %v попытки\n",2-count)
       auth:=authInput(count+1)
       if(auth!=nil){
         return auth
       }
       
     }
     return nil
}

