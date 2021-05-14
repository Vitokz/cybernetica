package model

type Groups struct {
	tableName         interface{} `sql:"groups"`
	Group_id          int64       `json:"group_id"`
	Parent_id         int64       `json:"parent_id"`
	Group_name        string      `json:"group_name"`
	Group_description string      `json:"group_description"`
}
