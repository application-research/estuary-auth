package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ApiKey struct {
	ApiKey string `json:"api_key"`
}

type UsernameApiKey struct {
	Username string `json:"api_key"`
	ApiKey   string `json:"api_key"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"passwordHash"`
}

func NewAuthRouterConfig() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//	initial connection to database

	fmt.Print("routes")
	// Routes
	e.POST("/check-api-key", BasicUserApiCheckHandler)
	e.POST("/check-user-api-key", BasicApiUserCheckHandler)
	e.POST("/check-user-pass", BasicUserPassHandler)
	e.GET("/register-new-token", BasicRegisterUserHandler)
	e.GET("/register-new-exp-token", BasicRegisterExpiringUserHandler)

	//	tricky redirection using a proxy
	// Start server
	fmt.Print("start")
	e.Logger.Fatal(e.Start("0.0.0.0:1313"))
}
