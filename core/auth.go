package core

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

// Authorization structures
type AuthorizationServer struct {
	// The authorization server's identifier.
	Authorization
}

type Authorization struct {
	DB     *gorm.DB
	tracer trace.Tracer
}

// Data models
type AuthToken struct {
	gorm.Model
	Token      string `gorm:"unique;->"` // read only to prevent storing new tokens but not break existing tokens
	TokenHash  string `gorm:"unique"`
	Label      string
	User       uint
	UploadOnly bool
	Expiry     time.Time
	IsSession  bool
}

type User struct {
	gorm.Model
	UUID     string `gorm:"unique"`
	Username string `gorm:"unique"`
	Salt     string
	PassHash string
	DID      string

	UserEmail string

	AuthToken AuthToken `gorm:"-"`
	Perm      int
	Flags     int

	StorageDisabled bool
}

// Initialize
func Init() *AuthorizationServer {
	return &AuthorizationServer{} // create the authorization server
}

// Sets a database connection.
func (s *AuthorizationServer) SetDB(db *gorm.DB) *AuthorizationServer {
	s.DB = db // connect to the database
	return s
}

// Set database connection with a string dsn
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

// Connect to the server and return the Authorization object
func (s *AuthorizationServer) Connect() Authorization {
	return s.Authorization
}

const TokenExpiryDurationPermanent = time.Hour * 24 * 365 * 100

func (s Authorization) NewUserAndAuthToken(expDuration time.Duration) (*AuthToken, error) {

	var user User
	s.DB.Model(&user).Where("id = ?", "260").First(&user)

	var expiry time.Time
	if expDuration == 0 {
		expiry = time.Now().Add(TokenExpiryDurationPermanent)
	} else {
		expiry = time.Now().Add(expDuration)
	}

	token := "EST" + uuid.New().String() + "ARY"
	fmt.Printf("Token: %s", token)
	authToken := &AuthToken{
		Token:      token,
		TokenHash:  GetTokenHash(token),
		Label:      "From DELTA",
		User:       user.ID,
		Expiry:     expiry,
		UploadOnly: true,
		IsSession:  true,
	}
	if err := s.DB.Create(authToken).Error; err != nil {
		return nil, err
	}

	return authToken, nil
}

// CheckAuthorizationToken Checking if the token is valid.
func (s Authorization) CheckAuthorizationToken(token string) (*User, error) {
	//cached, ok := s.cacher.Get(token)
	var authToken AuthToken
	tokenHash := GetTokenHash(token)
	if err := s.DB.First(&authToken, "token = ? OR token_hash = ?", token, tokenHash).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &HttpError{
				Code:    http.StatusUnauthorized,
				Reason:  ERR_INVALID_TOKEN,
				Details: "api key does not exist",
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

	user.AuthToken = authToken
	return &user, nil
}

// A middleware that checks if the user is authorized to access the API.
func (s Authorization) AuthRequired(level int) echo.MiddlewareFunc {
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

type ApiKeyParam struct {
	Username string
	Token    string `json:"api`
}
type AuthenticationResult struct {
	Username string     `json:"username,omitempty"`
	Password string     `json:"password,omitempty"`
	Salt     string     `json:"salt,omitempty"`
	Result   AuthResult `json:"result"`
}

type AuthResult struct {
	Validated bool   `json:"validated"`
	Details   string `json:"details"`
}
type AuthenticationParam struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s Authorization) AuthenticateApiKey(param ApiKeyParam) AuthenticationResult {

	var authToken AuthToken
	tokenHash := GetTokenHash(param.Token)
	if err := s.DB.First(&authToken, "token = ? OR token_hash = ?", param.Token, tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthenticationResult{
				Username: param.Username,
				Result: AuthResult{
					Validated: false,
					Details:   "api key does not exists",
				},
			}
		}
	}

	if authToken.Expiry.Before(time.Now()) {
		return AuthenticationResult{
			Username: param.Username,
			Result: AuthResult{
				Validated: false,
				Details:   "api key expired",
			},
		}
	}

	return AuthenticationResult{
		Username: param.Username,
		Result: AuthResult{
			Validated: true,
			Details:   "api key validated",
		},
	}
}

func (s Authorization) AuthenticateApiKeyUser(param ApiKeyParam) AuthenticationResult {
	var authToken AuthToken
	if err := s.DB.First(&authToken, "token = ?", param.Token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthenticationResult{
				Username: param.Username,
				Result: AuthResult{
					Validated: false,
					Details:   "api key does not exists",
				},
			}
		}
	}

	if authToken.Expiry.Before(time.Now()) {
		return AuthenticationResult{
			Username: param.Username,
			Result: AuthResult{
				Validated: false,
				Details:   ERR_TOKEN_EXPIRED,
			},
		}
	}

	var user User
	if err := s.DB.First(&user, "id = ?", authToken.User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthenticationResult{
				Username: param.Username,
				Result: AuthResult{
					Validated: false,
					Details:   "no user exists for the specified api key",
				},
			}
		}
	}

	return AuthenticationResult{
		Username: param.Username,
		Result: AuthResult{
			Validated: true,
			Details:   "api key and user is validated",
		},
	}
}
func (s Authorization) AuthenticateUserPassword(param AuthenticationParam) AuthenticationResult {

	var user User
	if err := s.DB.First(&user, "username = ?", strings.ToLower(param.Username)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthenticationResult{
				Username: param.Username,
				Result: AuthResult{
					Validated: false,
					Details:   "user not found",
				},
			}
		}
	}

	//	validate password
	//	SQLlite and Postgres has incompatibility in hashing and even though we are dropping support for sqlite later,
	//	we still need to accommodate those who chooses to use SQLite for experimentation purposes.
	var valid = true
	var dbDialect = s.DB.Config.Dialector.Name()

	//	check password hash (this is the way).
	if (user.Salt != "" && (user.PassHash != GetPasswordHash(param.Password, user.Salt, dbDialect))) || (user.Salt == "" && user.PassHash != param.Password) {
		valid = false                                                                            //	assume it's not valid.
		if bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(param.Password)) == nil { //	we are using bcrypt, so we need to rehash it.
			valid = true
		}
	}

	if !valid {
		return AuthenticationResult{
			Username: param.Username,
			Result: AuthResult{
				Validated: false,
				Details:   "user not found",
			},
		}
	}

	// Successful authentication
	return AuthenticationResult{
		Username: param.Username,
		Result: AuthResult{
			Validated: true,
			Details:   "user authenticated",
		},
	}

}
