package repository

import (
	"context"
	"main/model"
	"time"

	"github.com/sirupsen/logrus"
)

func (p *Pg) Auth(ctx context.Context, user *model.User) (*model.User, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, user) //Задаю таблицу и модель для запроса

	_, err := query.Insert() //Создаю запрос
	if err != nil {
		logrus.WithError(err).Error("Failed to insert in database")
		return user, err
	}
	user.Password = ""
	return user, err
}

func (p *Pg) Login(ctx context.Context, user *model.Login) (*model.Login, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()
	//Создаю запрос в таблицу с определенными условиями
	query := p.Db.ModelContext(cont, user).Where("login = ?", user.Login).Where("password = ?", user.Password)

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find this user in database")
		return user, err
	}
	return user, err
}
