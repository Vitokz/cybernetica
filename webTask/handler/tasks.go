package handler

import (
	"context"
	"main/model"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) TaskNew(c echo.Context) error {
	//readyTask := new(model.Task)
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "create new task",
	})
	task := new(model.Task)

	err := c.Bind(&task)
	if err != nil {
		log.WithError(err).Error("Incorrectly entered TaskNew json")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = h.checkGroup(ctx, task.Group_id)
	if err != nil {
		log.WithError(err).Error("group does not exist")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	task.Completed = false
	task.Created_at = time.Now()
	task.Task_id = hash(strconv.Itoa(task.Group_id) + task.Task)[:5]

	task, err = h.Db.TaskNew(ctx, task)
	if err != nil {
		log.WithError(err).Error("Failed to createing new task")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Println()
	return c.JSON(http.StatusOK, task)
}

func (h *Handler) TaskRefresh(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "refresh new task",
	})
	task := new(model.Task)
	id := c.Param("id")
	task.Task_id = c.Param("id")
	task, err = h.Db.InitTask(ctx, task)
	if err != nil {
		log.WithError(err).Error("Incorrect id")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = c.Bind(&task)
	if err != nil {
		log.WithError(err).Error("Incorrectly entered TaskNew json")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = h.checkGroup(ctx, task.Group_id)
	if err != nil {
		log.WithError(err).Error("group does not exist")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	task.Completed = false
	task.Task_id = hash(strconv.Itoa(task.Group_id) + task.Task)[:5]

	task, err = h.Db.TaskRefresh(ctx, task, id)
	if err != nil {
		log.WithError(err).Error("Failed to createing new task")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Println()
	return c.JSON(http.StatusOK, task)
}

func (h *Handler) TaskReady(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "Get ready or not ready task",
	})

	task := new(model.Task)
	task.Task_id = c.Param("id")
	task, err = h.Db.InitTask(ctx, task)
	if err != nil {
		log.WithError(err).Error("Incorrect id")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	status := c.QueryParam("finished")
	if status != "true" && status != "false" {
		log.Error("Get parametr is not bool")
		return echo.NewHTTPError(http.StatusInternalServerError)
	} else {
		if status == "true" {
			task.Completed = true
			task.Completed_at = time.Now()
		} else {
			task.Completed = false
			task.Completed_at = time.Now()
		}
	}

	task, err = h.Db.TaskReady(ctx, task)
	if err != nil {
		log.Error("Failed to updating completed")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Println()
	return c.JSON(http.StatusOK, task)
}

func (h *Handler) TaskSort(c echo.Context) error {
	var err error
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "Get time stat",
	})

	tasks := new([]model.Task)

	typeSort := c.QueryParam("sort")
	lm := c.QueryParam("limit")
	if lm == "" {
		lm = "10000000"
	}
	limit, err := strconv.Atoi(lm)
	if err != nil {
		log.Error("limit param is not int")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	typeSortCompleted := c.QueryParam("type")
	if typeSortCompleted != "all" && typeSortCompleted != "completed" && typeSortCompleted != "working" {
		log.Error("Get parametr type is wrong")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	switch typeSort {
	case "name":
		tasks, err := h.Db.TaskSortByName(ctx, tasks, limit, typeSortCompleted)
		if err != nil {
			log.Error("Failed to sort by name")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		log.Println()
		return c.JSON(http.StatusOK, tasks)
	case "groups":
		tasks, err := h.Db.TaskSortByGroups(ctx, tasks, limit, typeSortCompleted)
		if err != nil {
			log.Error("Failed to sort by groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		log.Println()
		return c.JSON(http.StatusOK, tasks) //Сделать праивильный вывод
	case "":
		tasks, err := h.Db.TaskSortByGroups(ctx, tasks, limit, typeSortCompleted)
		if err != nil {
			log.Error("Failed to sort by groups")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		log.Println()
		return c.JSON(http.StatusOK, tasks) //Сделать праивильный вывод
	default:
		log.Error("type parametr wrong")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
}

func (h *Handler) TaskStat(c echo.Context) error {
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "Task stat",
	})
	tasks := new(model.Task)
	typeSort := c.Param("type")

	result, err := h.Db.TaskStat(ctx, tasks, typeSort)
	if err != nil {
		log.WithError(err)
		return err
	}

	log.Println()
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) checkGroup(ctx context.Context, id int) error {
	group := new(model.Groups)
	err := h.Db.CheckGroup(ctx, group, id)
	if err != nil {
		return err
	}
	return nil
}

// func (h *Handler) checkId(ctx context.Context, id string) error {
// 	task := new(model.Task)
// 	err := h.Db.CheckId(ctx, task, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
