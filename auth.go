package auth

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AuthorizationServer struct {
	// The authorization server's identifier.
	DB     *gorm.DB
	Server interface{}
	tracer trace.Tracer
}

type AuthToken struct {
	gorm.Model
	Token      string `gorm:"unique"`
	User       uint
	UploadOnly bool
	Expiry     time.Time
}

type User struct {
	gorm.Model
	UUID     string `gorm:"unique"`
	Username string `gorm:"unique"`
	PassHash string
	DID      string

	UserEmail string

	authToken AuthToken
	Perm      int
	Flags     int

	StorageDisabled bool
}

//	Initialize
func Init() *AuthorizationServer {
	return &AuthorizationServer{}
}

//	Set DB
func (s *AuthorizationServer) setDB(db *gorm.DB) *AuthorizationServer {
	s.DB = db
	return s
}

func (s *AuthorizationServer) checkTokenAuth(token string) (*User, error) {
	var authToken AuthToken
	if err := s.DB.First(&authToken, "token = ?", token).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
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
	if err := s.DB.First(&user, "id = ?", authToken.User).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &HttpError{
				Code:    http.StatusUnauthorized,
				Reason:  ERR_INVALID_TOKEN,
				Details: "no user exists for the spicified api key",
			}
		}
		return nil, err
	}

	user.authToken = authToken
	return &user, nil
}

func (s *AuthorizationServer) AuthRequired(level int) echo.MiddlewareFunc {
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

			u, err := s.checkTokenAuth(auth)
			if err != nil {
				return err
			}

			span.SetAttributes(attribute.Int("user", int(u.ID)))

			if u.authToken.UploadOnly && level >= PermLevelUser {
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
