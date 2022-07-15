package auth

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//	Authorization structures
type AuthorizationServer struct {
	// The authorization server's identifier.
	*Authorization
}

type Authorization struct {
	DB     *gorm.DB
	tracer trace.Tracer
}

// Data models
type Token struct {
	gorm.Model
	Token      string `gorm:"unique"`
	User       uint
	UploadOnly bool
	Expiry     time.Time
}

type User struct {
	gorm.Model
	UUID            string `gorm:"unique"`
	Username        string `gorm:"unique"`
	PassHash        string
	DID             string
	UserEmail       string
	AuthToken       Token
	Perm            int
	Flags           int
	StorageDisabled bool
}

//	Initialize
func Init() *AuthorizationServer {
	return &AuthorizationServer{} // create the authorization server
}

//	Sets a database connection.
func (s *AuthorizationServer) SetDB(db *gorm.DB) *AuthorizationServer {
	s.DB = db // connect to the database
	return s
}

//	Set database connection with a string dsn
func (s *AuthorizationServer) SetDBWithString(dbConnection string) *AuthorizationServer {

	db, err := gorm.Open(postgres.Open(dbConnection), &gorm.Config{})
	if err != nil {
		panic(err) // database connection is required
	}

	s.DB = db // connect to the database
	return s
}

func (s *AuthorizationServer) SetDBConfig(dbConnection postgres.Config) *AuthorizationServer {

	db, err := gorm.Open(postgres.New(dbConnection), &gorm.Config{})

	if err != nil {
		panic(err) // database connection is required
	}

	s.DB = db // connect to the database
	return s
}

//	Connect to the server and return the Authorization object
func (s *AuthorizationServer) Connect() *Authorization {
	return s.Authorization
}

// Checking if the token is valid.
func (s *Authorization) CheckAuthorizationToken(token string, permission int) (*User, error) {
	var authToken Token
	if err := s.DB.First(&authToken, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &HttpError{
				Code:    http.StatusUnauthorized,
				Reason:  ERR_INVALID_TOKEN,
				Details: "api key does not exists",
			}
		}
		return nil, err
	}

	if authToken.Expiry.Before(time.Now()) {
		return nil, &HttpError{
			Code:    http.StatusUnauthorized,
			Reason:  ERR_TOKEN_EXPIRED,
			Details: fmt.Sprintf("token for user %d expired %s", authToken.User, authToken.Expiry),
		}
	}

	var user User
	if err := s.DB.First(&user, "id = ? and perm = ?", authToken.User, permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &HttpError{
				Code:    http.StatusUnauthorized,
				Reason:  ERR_INVALID_TOKEN,
				Details: "no user exists for the specified api key",
			}
		}
		return nil, err
	}

	user.AuthToken = authToken
	return &user, nil
}

// A middleware that checks if the user is authorized to access the API.
func (s *Authorization) AuthRequired(level int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//	Check first if the Token is available. We should not continue if the
			//	token isn't even available.
			auth, err := ExtractAuth(c)
			if err != nil {
				return err
			}

			ctx, span := s.tracer.Start(c.Request().Context(), "authCheck")
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))

			u, err := s.CheckAuthorizationToken(auth)
			if err != nil {
				return err
			}

			span.SetAttributes(attribute.Int("user", int(u.ID)))

			if u.AuthToken.UploadOnly && level >= PermLevelUser {
				return &HttpError{
					Code:    http.StatusForbidden,
					Reason:  ERR_NOT_AUTHORIZED,
					Details: "api key is upload only",
				}
			}

			if u.Perm >= level {
				c.Set("user", u)
				return next(c)
			}

			return &HttpError{
				Code:    http.StatusForbidden,
				Reason:  ERR_NOT_AUTHORIZED,
				Details: "user not authorized",
			}
		}
	}
}
