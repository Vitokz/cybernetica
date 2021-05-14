package middlewares

import (
	"main/proto"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CheckJwtAuth(g *echo.Group) { //это middleware для проверки токена при переходе на другие эндпоинты
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(proto.JWTKEY),
		TokenLookup:   "cookie:JWTCookie",
	}))
}
