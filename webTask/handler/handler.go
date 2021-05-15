package handler

import (
	"context"
	"main/model"
)

type Handler struct {
	Db Repository
}

type Repository interface {
	AuthRepository
	GroupRepository
	TaskRepository
}

type AuthRepository interface {
	Auth(context.Context, *model.User) (*model.User, error)
	Login(context.Context, *model.Login) (*model.Login, error)
}

type GroupRepository interface {
	Groups(context.Context, *[]model.Groups) (*[]model.Groups, error)
	// GroupSort()
	// GroupTop()
	GroupId(context.Context, *model.GroupNew, int) (*model.GroupNew, error)
	// GroupChildsByID()
	GroupNew(context.Context, *model.GroupNew) (*model.GroupNew, error)
	// GroupRefresh()
	// GroupDelete()
}

type TaskRepository interface {
	// TaskNew()
	// TaskRefresh()
	// TaskStat()
	// TaskReady()
	// TaskSort()
}
