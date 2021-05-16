package model

type Groups struct {
	tableName         interface{} `sql:"groups"`
	Group_id          int64       `json:"group_id"`
	Parent_id         int64       `json:"parent_id"`
	Group_name        string      `json:"group_name"`
	Group_description string      `json:"group_description"`
}

type GroupsSort struct { //Незаконченнаое
	tableName interface{} `sql:"groups"`
}

type GroupNew struct {
	tableName         interface{} `sql:"groups"`
	Parent_id         int64       `json:"parent_id"`
	Group_name        string      `json:"group_name"`
	Group_description string      `json:"group_description"`
}

type GroupParents struct {
	tableName         interface{} `sql:"groups"`
	Group_id          int64       `json:"group_id"`
	Parent_id         int64       `json:"parent_id"`
	Group_name        string      `json:"group_name"`
	Group_description string      `json:"group_description"`
}
