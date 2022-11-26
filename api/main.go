package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "net/http"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"passwordHash"`
}

func NewAuthRouterConfig() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//	initial connection to database

	// Routes
	e.GET("/check-api-key", BasicApiCheckHandler)
	e.GET("/check-user-passhash", BasicUserPassHashHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1000"))
}

// Handler
func BasicApiCheckHandler(c echo.Context) error {
	//authorizationServer := new(core.AuthorizationServer)
	//auth := authorizationServer.SetDB()
}

// Handler
func BasicUserPassHashHandler(c echo.Context) error {
	//c.Request().Header.Set("Authorization
}
