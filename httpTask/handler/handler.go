package handler

import (
  "net/http"
  "encoding/json"
  "crypto/sha256"
  "fmt"
  "main/models"
  "main/proto"
  "main/repasitory"
)

type Handler struct{
  DB *repasitory.Db
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json") //Задаю тип данных с которыми работаю

  var newUser models.User //Создаю локальную переменную Юзера
  _ = json.NewDecoder(r.Body).Decode(&newUser) // Связываем тело REQUEST с нашей переменной нового юзера

  h.DB.Users=append(h.DB.Users,newUser) // Добавляю новоиспченного пользователя в "базу данных"
  json.NewEncoder(w).Encode(newUser) //вывожу данные о новом пользоватле
}


func(h *Handler) AdminUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")// Задаю тип данных json
  var sess models.UserSession //Создаю локальную струкутру с токеном и ролью пользователя
  _=json.NewDecoder(r.Body).Decode(&sess) // связываю ntkj REQUEST с переменной sess в которую я передавал token польователя 
  for _,v := range h.DB.Session{ // Цикл для проверки наличия пользователя в базе и поиск его роли в этом имре
    if(v.Token==sess.Token){
      sess.Role=v.Role
    }
  }
 
  if(sess.Role=="") { // Проверка на авторизацию
    w.WriteHeader(401) // выдаю ошибку 401
    json.NewEncoder(w).Encode("Вы не авторизованы")
    return
  }
  if(sess.Role!="admin") { // Проверка на админа
    w.WriteHeader(400) // выдаю ошибку 400
    json.NewEncoder(w).Encode("Вы не админ")
    return
  }

  answer := proto.Address{
    AddressInfo: "adad", // Вывожу определенный JSon в случае успеха проверки
  }

  json.NewEncoder(w).Encode(answer) 
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json") //Задаю тип данных с которыми работаю
  
  var auth proto.LoginSchema  // Переменная содержащаяя переданные данные о пользователе(логин и пароль)
  _=json.NewDecoder(r.Body).Decode(&auth) // связываю переменную с REQUEST

  login:=auth.Login  // Достаю логин
  password:=auth.Password //Достаю пароль

  if(login== "" || password== "") { //Проверка на пустоту полей авторизации
    err := fmt.Errorf("Поля логина или пароля пустуют") 
    fmt.Println(err)
  }

  var sess models.UserSession // Переменная со структурой models.UserSession (содержащая token и роль)
  how:=sha256.New() //Класс хэширования

  check:=false // Проверка нашелся ли юзер
  for _,v:=range h.DB.Users { //Перебираю весь массив на поиск нужного юзера 
    if(v.Login==login && v.Password==password){
      how.Write([]byte(v.Name)) //Добаляю в хэш буфер имя
      how.Write([]byte(login)) //Добавлю ы хэш буфер логин
      how.Write([]byte(password)) // Добавлю в хэш буфер логин
      how.Write([]byte(v.Role)) // Добавлю в хэш буфер роль
      sess.Role=v.Role // Добаляю роль в глоб переменную
      check=true 
      break
    }
  }

  if(check==false) {
    w.WriteHeader(400) // Возвращаю ошибку 400 если юзер не обнаружен
  }else{
    sess.Token=fmt.Sprintf("%x",how.Sum(nil)) //достаю из хэш буфера свой токен
    h.DB.Session=append(h.DB.Session,sess) // добавлю хэш пользователя в глоб переменную
    json.NewEncoder(w).Encode(sess) // вывожу токен залогиневшегос япользователя  
    json.NewEncoder(w).Encode(h.DB.Users) //вывожу всех юзеров
    json.NewEncoder(w).Encode(h.DB.Session) // вывожу все токены
  }
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json") //ЗАдаю тип данных с которыми буду работать
  js:=make(map[string]string) //Создаю ассоциативный массив который в будущем будет содержать json переданный фаил
  _=json.NewDecoder(r.Body).Decode(&js) // связываю переменную js с переданным файлом
  
  check:=false //переменная для проверки 
  for _,v := range h.DB.Session{//Проверяю есть ли такой пользователь в базе с токенами и достаю его роль в случае ууспеха
    if(v.Token==js["token"]){
     check=true
    }
  }
  if(check==false) { // если он не найден вывожу ошибку 401
    w.WriteHeader(401)
  }

  json.NewEncoder(w).Encode(js) //вывожу переданный поьзователем json
}
