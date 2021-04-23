package models

type ( // структура для хранения данных о пользователе
  User struct {
    Name     string `json:"name"`
    Login    string `json:"login"`
    Password string `json:"password"`
    Role     string `json:"role"`
  }
)