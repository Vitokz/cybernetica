package main

import (
	"main/handler"
	"main/repository"
	"main/server"

	"github.com/labstack/echo/v4"
)

//"main/handler"
//"net/http"
//"github.com/labstack/echo/middleware"
//"github.com/labstack/echo/v4"
//"github.com/labstack/echo/v4/middleware"

func main() {
	pgClient := repository.NewPgClient()
	db := pgClient.GetConnection()
	defer db.Close()

	hndlr := handler.Handler{
		Db: &repository.Pg{Db: db},
	}

	router := echo.New()

	rest := server.Rest{
		Handler: &hndlr,
		Router:  router,
	}

	rest.Rout()

	rest.Start(":8080")

	// e := echo.New()

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// //Auth
	// e.POST("/auth/new", handler.Auth)
	// e.POST("/login", handler.Login)

	// //Groups
	// e.GET("/groups", handler.Groups)
	// e.GET("/groups?sort=name&limit=n", handler.GroupsSort)
	// e.GET("/group/top_parents", handler.GroupTop)
	// e.GET("/group/:id", handler.GroupId)
	// e.GET("/group/childs/:id", handler.GroupChildsByID)
	// e.POST("/group/new", handler.GroupNew)
	// e.PUT("/group/:id", handler.GroupRefresh)
	// e.DELETE("/group/:id", handler.GroupDelete)

	// //Tasks
	// e.GET("/tasks?sort=[name|groups]&limit=n&type=[all|completed|working]", handler.TasksSort)
	// e.POST("/tasks/new", handler.TaskNew)
	// e.PUT("/tasks/:id", handler.TaskRefresh)
	// e.POST("/tasks/:id?finished=(true|false)", handler.TaskReady)
	// e.GET("/stat/(today|yesterday|week|month)", handler.TaskStat)
	// e.Logger.Fatal(e.Start(":1323"))
}
