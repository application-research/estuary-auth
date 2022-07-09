package auth

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestAuthorizationServer_AuthRequired(t *testing.T) {
	type fields struct {
		DB     *gorm.DB
		Server interface{}
		tracer trace.Tracer
	}
	type args struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   echo.MiddlewareFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthorizationServer{
				DB:     tt.fields.DB,
				Server: tt.fields.Server,
				tracer: tt.fields.tracer,
			}
			if got := s.AuthRequired(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationServer_checkTokenAuth(t *testing.T) {
	type fields struct {
		DB     *gorm.DB
		Server interface{}
		tracer trace.Tracer
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthorizationServer{
				DB:     tt.fields.DB,
				Server: tt.fields.Server,
				tracer: tt.fields.tracer,
			}
			got, err := s.checkTokenAuth(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkTokenAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkTokenAuth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationServer_setDB(t *testing.T) {
	type fields struct {
		DB     *gorm.DB
		Server interface{}
		tracer trace.Tracer
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *AuthorizationServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthorizationServer{
				DB:     tt.fields.DB,
				Server: tt.fields.Server,
				tracer: tt.fields.tracer,
			}
			if got := s.setDB(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		want *AuthorizationServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
