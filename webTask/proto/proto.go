package proto

import "main/model"

const (
	JWTKEY = "nelez" //Секретный ключ для токена
	SALT   = "asd"   //Соль для хэша
)

type LoginResponse struct {
	Token string `json:"token"`
}

type GroupsResponse struct {
	Groups []model.Groups `json:"groups"`
}
