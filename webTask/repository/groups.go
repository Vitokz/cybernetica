package repository

import (
	"context"
	"main/model"
	"time"

	"github.com/go-pg/pg/orm"
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

func (p *Pg) GroupChildsByID(ctx context.Context, groups *[]model.GroupNew, id int) (*[]model.GroupNew, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Where("parent_id = ?", id)

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find groups")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupTop(ctx context.Context, groups *[]model.GroupParents) (*[]model.GroupParents, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Where("parent_id = ?", 0)

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find groups")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupRefresh(ctx context.Context, groups *model.GroupNew, id int) (*model.GroupNew, error) {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Set("parent_id = ?", groups.Parent_id).Set("group_name = ?", groups.Group_name).Set("group_description = ?", groups.Group_description).Where("group_id = ?", id)

	_, err := query.Update() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No Updated group")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupDelete(ctx context.Context, groups *model.Groups, id int) error {
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	query := p.Db.ModelContext(cont, groups).Where("group_id = ?", id)

	_, err := query.Delete() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find groups")
		return err
	}

	return nil
}

func (p *Pg) GroupsSortByName(ctx context.Context, groups *[]model.Groups, limit int) (*[]model.Groups, error) {
	var query *orm.Query
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	if limit == -1 {
		query = p.Db.ModelContext(cont, groups).Order("group_name ASC")
	} else {
		query = p.Db.ModelContext(cont, groups).Order("group_name ASC").Limit(limit)
	}

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find groups")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupsSortByParentsFirst(ctx context.Context, groups *[]model.Groups, limit int) (*[]model.Groups, error) {
	var query *orm.Query
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	if limit == -1 {
		query = p.Db.ModelContext(cont, groups).Order("parent_id")
	} else {
		query = p.Db.ModelContext(cont, groups).Order("parent_id").Limit(limit)
	}

	err := query.Select() //Отправляю запрос
	if err != nil {
		logrus.WithError(err).Error("No find groups")
		return groups, err
	}

	return groups, nil
}

func (p *Pg) GroupsByParentsWithChilds(ctx context.Context, groups *[]model.Groups, limit int) (*[]model.Groups, error) {
	//var query *orm.Query
	var result []model.Groups
	cont, canceelFunc := context.WithTimeout(ctx, 30*time.Second) //Создаю ограничение по времени запроса
	defer canceelFunc()

	count := new(model.Groups)
	countQuery, err := p.Db.ModelContext(cont, count).Where("parent_id = ?", 0).Count()
	if err != nil {
		logrus.WithError(err).Error("No find parents")
		return groups, err
	} else if countQuery > limit && limit != -1 {
		countQuery = limit
	}

	parents := make([]model.Groups, 0)
	err = p.Db.ModelContext(cont, &parents).Where("parent_id = ?", 0).Order("parent_id").Select()
	parents = reverse(parents)

	for i := 0; i < countQuery; i++ {
		result = append(result, parents[i])

		childs := make([]model.Groups, 0)
		err = p.Db.ModelContext(cont, &childs).Where("parent_id = ?", parents[i].Group_id).Select()
		if err != nil {
			continue
		}

		result = append(result, childs[:]...)
	}
	groups = &result

	return groups, nil
}

func reverse(input []model.Groups) []model.Groups {
	if len(input) == 0 {
		return input
	}
	return append(reverse(input[1:]), input[0])
}
