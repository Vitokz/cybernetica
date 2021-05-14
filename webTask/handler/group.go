package handler

import (
	"main/model"
	"main/proto"
	"net/http"

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
	}

	return c.JSON(http.StatusOK, proto.GroupsResponse{
		Groups: *groups,
	})
}

// func (h *Handler) Groups
