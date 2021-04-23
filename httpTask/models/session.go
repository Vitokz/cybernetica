package models 


type ( //стркутура для хранения токенов и роли пользоватлей
    UserSession struct {
    Token string `json:"token"`
    Role string `json:"-"`
  }
)