package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "main/handler"
  "main/repasitory"
  "main/models"
)


func main() {

db := repasitory.Db{   //Создаю структуру DB которая лежит в папке repasitory 
  Users: []models.User{}, // Добавляю массив типа models.User{} для работы с юзерами
  Session: []models.UserSession{}, // Добавляю массив типа models.UserSession{} для работы с токенами юзеров
}
 hndlr := handler.Handler{ // Создаю тип для добавления инструкций к каждому эндпоинту + добавляю туда нашу структуру для работы с массивами ползователей и токенов
   DB: &db,
 }

 r := mux.NewRouter() // Создаю роутер для перехода между эндпоинтами
 r.HandleFunc("/create", hndlr.CreateUser).Methods("POST") //Добавляю эндпоинт /create для создания пользвателя и использую инструкцию hndlr.CreateUser с методом POST
 r.HandleFunc("/login", hndlr.LoginUser).Methods("POST") //Добаляю эндпоинт /login для авторизации пользователя с использованием инструкции hndlr.LoginUser
 r.HandleFunc("/addres", hndlr.AdminUser).Methods("POST")// Добавляю эндпоинт / addres для проверки наадмина 
 r.HandleFunc("/info",hndlr.Info).Methods("POST")// ДОбавляю эндпоинт /info для вывода Json файла введеного пользоватлем 
 log.Fatal(http.ListenAndServe(":8000", r))
}

