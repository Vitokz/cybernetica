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
