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
	GroupTop(context.Context, *[]model.GroupParents) (*[]model.GroupParents, error)
	GroupId(context.Context, *model.GroupNew, int) (*model.GroupNew, error)
	GroupChildsByID(context.Context, *[]model.GroupNew, int) (*[]model.GroupNew, error)
	GroupNew(context.Context, *model.GroupNew) (*model.GroupNew, error)
	GroupRefresh(context.Context, *model.GroupNew, int) (*model.GroupNew, error)
	GroupDelete(context.Context, *model.Groups, int) error
	//GroupSort functions
	GroupsSortByName(context.Context, *[]model.Groups, int) (*[]model.Groups, error)
	GroupsSortByParentsFirst(context.Context, *[]model.Groups, int) (*[]model.Groups, error)
	GroupsByParentsWithChilds(context.Context, *[]model.Groups, int) (*[]model.Groups, error)
}

type TaskRepository interface {
	TaskNew(context.Context, *model.Task) (*model.Task, error)
	TaskRefresh(context.Context, *model.Task, string) (*model.Task, error)
	TaskStat(context.Context, *model.Task, string) (*model.TaskStat, error)
	TaskReady(context.Context, *model.Task) (*model.Task, error)
	// TaskSort
	TaskSortByName(context.Context, *[]model.Task, int, string) (*[]model.Task, error)
	TaskSortByGroups(context.Context, *[]model.Task, int, string) (*[]model.Task, error)
	//Вспомогательные
	CheckGroup(context.Context, *model.Groups, int) error
	InitTask(context.Context, *model.Task) (*model.Task, error)
}
