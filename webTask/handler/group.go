package handler

import (
	"main/model"
	"main/proto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Groups(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "check all groups",
	})

	groups := new([]model.Groups)

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

// func (h *Handler) GroupsSort(c echo.Context) error {
// 	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
// 	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
// 		"event": "check all groups",
// 	})
// 	typeSort := c.QueryParam("sort")
// 	limit := c.QueryParam("limit")

// 	switch typeSort {
// 	case "name": //Сортировка по имени

// 	case "parents_first": //Сначала идут все родители птом все дети, Все отсортированы по алфавиту

// 	case "parent_with_childs": //Сначала идет родитель потом его дети и тд

// 	}
// }

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
		"event": "insert new group",
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
	return c.JSON(http.StatusOK, groups)
}
