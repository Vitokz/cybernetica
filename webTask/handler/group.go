package handler

import (
	"main/model"
	"main/proto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GroupsSort(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "sort all groups",
	})

	typeSort := c.QueryParam("sort")
	lm := c.QueryParam("limit")
	if lm == "" {
		lm = "-1"
	}

	limit, err := strconv.Atoi(lm)
	if err != nil {
		log.Error("limit param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	groups := new([]model.Groups)

	switch typeSort {
	case "name": //Сортировка по имени

		log.WithFields(logrus.Fields{
			"sort by": "name",
		})

		groups, err = h.Db.GroupsSortByName(ctx, groups, limit)
		if err != nil {
			log.WithError(err).Error("Failed to sorting groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		log.Println()
		return c.JSON(http.StatusOK, proto.GroupsResponse{
			Groups: *groups,
		})

	case "parents_first": //Сначала идут все родители птом все дети, Все отсортированы по алфавиту
		log.WithFields(logrus.Fields{
			"sort by": "parents_firts",
		})

		groups, err = h.Db.GroupsSortByParentsFirst(ctx, groups, limit)
		if err != nil {
			log.WithError(err).Error("Failed to sorting groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		log.Println()
		return c.JSON(http.StatusOK, proto.GroupsResponse{
			Groups: *groups,
		})
	case "parent_with_childs": //Сначала идет родитель потом его дети и тд
		log.WithFields(logrus.Fields{
			"sort by": "parent_with_childs",
		})

		groups, err = h.Db.GroupsByParentsWithChilds(ctx, groups, limit)
		if err != nil {
			log.WithError(err).Error("Failed to select groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		log.Println()
		return c.JSON(http.StatusOK, proto.GroupsResponse{
			Groups: *groups,
		})
	case "":
		log.WithFields(logrus.Fields{
			"sort by": "random",
		})

		groups, err = h.Db.Groups(ctx, groups)
		if err != nil {
			log.WithError(err).Error("Failed to select groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		log.Println()
		return c.JSON(http.StatusOK, proto.GroupsResponse{
			Groups: *groups,
		})
	}
	return echo.NewHTTPError(http.StatusInternalServerError)
}

func (h *Handler) GroupNew(c echo.Context) error {
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "insert new group",
	})

	newGroup := new(model.GroupNew)

	err := c.Bind(&newGroup)
	if err != nil {
		log.WithError(err).Error("Failed to accept new group json")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	newGroup, err = h.Db.GroupNew(ctx, newGroup)
	if err != nil {
		log.WithError(err).Error("No inserting new group to DB")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.WithFields(logrus.Fields{
		"group_name": newGroup.Group_name,
	}).Println()
	return c.JSON(http.StatusOK, newGroup)
}

func (h *Handler) GroupId(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "select by id",
	})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Id param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	groups := new(model.GroupNew)

	groups, err = h.Db.GroupId(ctx, groups, id)
	if err != nil {
		log.WithError(err).Error("Id not exist")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, groups)
}

func (h *Handler) GroupChildsByID(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "selected childs by parent id",
	})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Id param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	groups := new([]model.GroupNew)

	groups, err = h.Db.GroupChildsByID(ctx, groups, id)
	if err != nil {
		log.WithError(err).Error("Id not exist")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, proto.GroupsChildsById{
		Ghilds: *groups,
	})
}

func (h *Handler) GroupTop(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "select all parents",
	})

	groups := new([]model.GroupParents)

	groups, err = h.Db.GroupTop(ctx, groups)
	if err != nil {
		log.WithError(err).Error("Failed parents select")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, proto.GroupsParents{
		Parents: *groups,
	})
}

func (h *Handler) GroupRefresh(c echo.Context) error {
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "refresh a group",
	})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Id param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	newGroup := new(model.GroupNew)

	err = c.Bind(&newGroup)
	if err != nil {
		log.WithError(err).Error("Failed to accept group json")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	newGroup, err = h.Db.GroupRefresh(ctx, newGroup, id)
	if err != nil {
		log.WithError(err).Error("No refresh group to DB")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, newGroup)
}

func (h *Handler) GroupDelete(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "delete group",
	})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Id param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	groups := new(model.Groups)

	err = h.Db.GroupDelete(ctx, groups, id)
	if err != nil {
		log.WithError(err).Error("Id not exist")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, "Group deleted")
}
