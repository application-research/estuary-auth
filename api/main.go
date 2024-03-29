package main

import (
	"fmt"
	"github.com/application-research/estuary-auth/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	_ "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (

	// RuntimeVer date string of when build was performed filled in by -X compile flag
	auth *core.AuthorizationServer

	// OsSignal signal used to shutdown
	OsSignal chan os.Signal
)

func main() {
	OsSignal = make(chan os.Signal, 1)
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	dbHost, okHost := viper.Get("DB_HOST").(string)
	dbUser, okUser := viper.Get("DB_USER").(string)
	dbPass, okPass := viper.Get("DB_PASS").(string)
	dbName, okName := viper.Get("DB_NAME").(string)
	dbPort, okPort := viper.Get("DB_PORT").(string)
	if !okHost || !okUser || !okPass || !okName || !okPort {
		panic("invalid database configuration")
	}

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=prefer TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sqldb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqldb.SetMaxIdleConns(80)
	sqldb.SetMaxOpenConns(250)
	sqldb.SetConnMaxIdleTime(time.Hour)
	sqldb.SetConnMaxLifetime(time.Second * 60)
	if err != nil {
		panic("failed to connect database")
	}
	auth = core.Init()
	auth = auth.SetDB(db)

	NewAuthRouterConfig()

	LoopForever()
}

// LoopForever on signal processing
func LoopForever() {
	fmt.Printf("Entering infinite loop\n")

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	_ = <-OsSignal

	fmt.Printf("Exiting infinite loop received OsSignal\n")
}

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

func BasicRegisterUserHandler(c echo.Context) error {
	result, _ := auth.NewUserAndAuthToken(0)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   result.Token,
		"expires": result.Expiry,
	})

}

// BasicRegisterExpiringUserHandler // such as "300ms", "-1.5h" or "2h45m".
// // Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func BasicRegisterExpiringUserHandler(c echo.Context) error {
	duration := c.QueryParam("duration")
	fmt.Print(duration)
	durationToParse, err := time.ParseDuration(duration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid duration",
		})
	}
	result, _ := auth.NewUserAndAuthToken(durationToParse)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   result.Token,
		"expires": result.Expiry,
	})

}

// BasicUserApiCheckHandler `BasicUserApiCheckHandler` is a function that takes a `echo.Context` and returns an `error`
func BasicUserApiCheckHandler(c echo.Context) error {
	var apiKeyParam core.ApiKeyParam
	if err := c.Bind(&apiKeyParam); err != nil {
		return err
	}
	result := auth.AuthenticateApiKey(apiKeyParam)
	return c.JSON(http.StatusOK, result)
}

// Handler
// `BasicApiUserCheckHandler` is a function that takes a `echo.Context` and returns an `error`
func BasicApiUserCheckHandler(c echo.Context) error {
	var apiKeyParam core.ApiKeyParam
	if err := c.Bind(&apiKeyParam); err != nil {
		return err
	}
	result := auth.AuthenticateApiKeyUser(apiKeyParam)
	return c.JSON(http.StatusOK, result)
}

// Handler
// `BasicUserPassHandler` is a function that takes a `echo.Context` and returns an `error`
func BasicUserPassHandler(c echo.Context) error {
	var authParam core.AuthenticationParam
	if err := c.Bind(&authParam); err != nil {
		return err
	}

	result := auth.AuthenticateUserPassword(authParam)
	return c.JSON(http.StatusOK, result)
}
