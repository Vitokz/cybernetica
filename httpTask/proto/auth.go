package proto

type( //Схема json файла который должен передаваться при авторизации
  LoginSchema struct{
    Login string `json:"login"`
    Password string `json:"password"`
  }
)