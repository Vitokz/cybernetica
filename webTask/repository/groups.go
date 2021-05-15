package repository

import (
	"context"
	"main/model"
	"time"

	"github.com/sirupsen/logrus"
)

func (p *Pg) Groups(ctx context.Context, groups *[]model.Groups) (*[]model.Groups, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups)

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No selected groups")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupNew(ctx context.Context, groups *model.GroupNew) (*model.GroupNew, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups)

	_, err := query.Insert() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No inserting New group")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupId(ctx context.Context, groups *model.GroupNew, id int) (*model.GroupNew, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Where("group_id = ?", id)

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No selected groups")
		return groups, err
	}

	return groups, nil
}
