package server

import (
	"main/handler"

	"github.com/labstack/echo/v4"
)

type Rest struct {
	Handler *handler.Handler
	Router  *echo.Echo
}

func (r *Rest) Rout() {
	g := r.Router.Group("/groups")
	g.GET("", r.Handler.Groups)
	//g.GET("?sort=name&limit=n", r.Handler.GroupsSort)
	// g.GET("/top_parents", r.Handler.GroupTop)
	g.GET("/:id", r.Handler.GroupId)
	// g.GET("/childs/:id", r.Handler.GroupChildsByID)
	g.POST("/new", r.Handler.GroupNew)
	// g.PUT("/:id", r.Handler.GroupRefresh)
	// g.DELETE("/:id", r.Handler.GroupDelete)

	r.Router.POST("/auth/new", r.Handler.Auth)
	r.Router.POST("/login", r.Handler.Login)

	// t := r.Router.Group("/tasks")
	// t.GET("?sort=[name|groups]&limit=n&type=[all|completed|working]", r.Handler.TasksSort)
	// t.POST("/new", r.Handler.TaskNew)
	// t.PUT("/:id", r.Handler.TaskRefresh)
	// t.POST("/:id?finished=(true|false)", r.Handler.TaskReady)
	// t.GET("/stat/(today|yesterday|week|month)", r.Handler.TaskStat)
}

func (r *Rest) Start(port string) {
	r.Router.Logger.Fatal(r.Router.Start(port))
}
