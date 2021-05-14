package model

type Login struct {
	tableName interface{} `sql:"users"`
	Id        int64       `json:"id"`
	Login     string      `json:"login"`
	Password  string      `json:"password"`
}
