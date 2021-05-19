package repository

import (
	"context"
	"errors"
	"main/model"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

func (p *Pg) TaskNew(ctx context.Context, task *model.Task) (*model.Task, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, task)

	_, err := query.Insert() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("Failed to add")
		return task, err
	}

	_, err = p.Db.ModelContext(cont, task).Set("completed = ?", false).Where("task_id = ?", task.Task_id).UpdateNotNull()

	return task, nil
}

func (p *Pg) CheckGroup(ctx context.Context, groups *model.Groups, id int) error {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Where("group_id = ?", id)

	err := query.Select()
	if err != nil {
		logrus.WithError(err).Error("Group with this id does not exist")
		return err
	}
	return nil
}

func (p *Pg) TaskRefresh(ctx context.Context, task *model.Task, id string) (*model.Task, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, task).Where("task_id = ?", id)

	_, err := query.UpdateNotNull() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh")
		return task, err
	}

	return task, nil
}

func (p *Pg) TaskReady(ctx context.Context, task *model.Task) (*model.Task, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, task).Set("completed = ?", task.Completed).Where("task_id = ?", task.Task_id)

	_, err := query.UpdateNotNull() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("Failed to update")
		return task, err
	}

	return task, nil
}

func (p *Pg) InitTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, task).Where("task_id= ?", task.Task_id)
	err := query.Select()
	if err != nil {
		logrus.WithError(err).Error("Task with this id does not exist")
		return nil, err
	}
	return task, nil
}

func (p *Pg) TaskSortByName(ctx context.Context, task *[]model.Task, limit int, typeSort string) (*[]model.Task, error) {
	var query *orm.Query
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	switch typeSort {
	case "all":
		query = p.Db.ModelContext(ctx, task).Limit(limit)
	case "completed":
		query = p.Db.ModelContext(cont, task).Where("completed = ?", true).Limit(limit)
	case "working":
		query = p.Db.ModelContext(cont, task).Where("completed = ?", false).Limit(limit)
	}
	err := query.Select()
	if err != nil {
		logrus.WithError(err).Error("Failed to query sorting")
		return nil, err
	}
	return task, nil
}

func (p *Pg) TaskSortByGroups(ctx context.Context, task *[]model.Task, limit int, typeSort string) (*[]model.Task, error) {
	var query *orm.Query
	ctx, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	switch typeSort {
	case "all":
		query = p.Db.ModelContext(ctx, task).Order("group_id").Limit(limit)
	case "completed":
		query = p.Db.ModelContext(ctx, task).Where("completed = ?", true).Order("group_id").Limit(limit)
	case "working":
		query = p.Db.ModelContext(ctx, task).Where("completed = ?", false).Order("group_id").Limit(limit)
	}
	err := query.Select()
	if err != nil {
		logrus.WithError(err).Error("Failed to query sorting")
		return nil, err
	}
	return task, nil
}

func (p *Pg) TaskStat(ctx context.Context, task *model.Task, typeSort string) (*model.TaskStat, error) {
	//var query *orm.Query
	var err error
	result := new(model.TaskStat)
	ctx, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	switch typeSort {
	case "today":
		result.Created, err = p.Db.ModelContext(ctx, task).Where("created_at > ?", time.Now().Add(-24*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
		result.Completed, err = p.Db.ModelContext(ctx, task).Where("completed_at > ?", time.Now().Add(-24*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
	case "yesterday":
		result.Created, err = p.Db.ModelContext(ctx, task).Where("created_at > ?", time.Now().Add(-48*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
		result.Completed, err = p.Db.ModelContext(ctx, task).Where("completed_at > ?", time.Now().Add(-48*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
	case "week":
		result.Created, err = p.Db.ModelContext(ctx, task).Where("created_at > ?", time.Now().Add(-168*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
		result.Completed, err = p.Db.ModelContext(ctx, task).Where("completed_at > ?", time.Now().Add(-168*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
	case "month":
		result.Created, err = p.Db.ModelContext(ctx, task).Where("created_at > ?", time.Now().Add(-720*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}
		result.Completed, err = p.Db.ModelContext(ctx, task).Where("completed_at > ?", time.Now().Add(-720*time.Hour)).Count()
		if err != nil {
			logrus.WithError(err).Error("Failed to calculate stat created : today")
			return nil, err
		}

	default:
		err := errors.New("Get parametr is wrong")
		return nil, err
	}
	return result, nil
}

func deleteDuplicates(arr []int) []int {
	result := make([]int, 0)
	for _, value := range arr {
		if contains(result, value) {
			continue
		}
		result = append(result, value)
	}
	return result
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
