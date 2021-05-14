package model

type User struct {
	tableName interface{} `sql:"users"`
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	Login     string      `json:"login"`
	Password  string      `json:"password"`
}
