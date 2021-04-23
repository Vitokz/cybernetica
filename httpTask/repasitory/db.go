package repasitory

import "main/models"

type ( //СТруктура нашей имитированной базы данных с массивами пользователей и их ролей с токенами
  Db struct{
   Users []models.User
   Session []models.UserSession
  }
)