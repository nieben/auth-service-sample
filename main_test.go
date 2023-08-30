package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/api"
	"github.com/nieben/auth-service-sample/model"
	"github.com/nieben/auth-service-sample/route"
	"github.com/nieben/auth-service-sample/route/middleware"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	router *gin.Engine

	initToken      = ""
	expireToken    = ""
	bobToken       = &initToken
	bobExpireToken = &expireToken // test for expire
	invalidToken   *string
)

func init() {
	port = 8080
	model.TokenLifeTime = 5

	gin.SetMode(gin.TestMode)

	router = route.Init()

	s := "12345612345612345612345612345612"
	invalidToken = &s
}

func post(uri, method string, param interface{}, headers map[string]*string, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(method, uri, bytes.NewReader(jsonByte))
	for k, v := range headers {
		req.Header.Set(k, *v)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// test in one flow for all endpoints and cases
// user: bob eve
// role: admin ops dev
func TestFullFlow(t *testing.T) {
	cases := []struct {
		path       string
		method     string
		name       string
		expectCode int64
		expectErr  string
		param      interface{}
		header     map[string]*string
	}{
		// user operation
		{
			path:       "/user/create",
			method:     "POST",
			name:       "name err",
			expectCode: 1,
			expectErr:  model.UserNameErr.Error(),
			param: api.CreateUser{
				Username: "bob~^&",
				Password: "123456",
			},
		},
		{
			path:       "/user/create",
			method:     "POST",
			name:       "pwd err",
			expectCode: 1,
			expectErr:  model.UserPwdErr.Error(),
			param: api.CreateUser{
				Username: "bob",
				Password: "123",
			},
		},
		{
			path:       "/user/create",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.CreateUser{
				Username: "bob",
				Password: "123456",
			},
		},
		{
			path:       "/user/create",
			method:     "POST",
			name:       "user exist",
			expectCode: 1,
			expectErr:  model.UserExistErr.Error(),
			param: api.CreateUser{
				Username: "bob",
				Password: "123456",
			},
		},
		{
			path:       "/user/create",
			method:     "POST",
			name:       "ok 2",
			expectCode: 0,
			expectErr:  "",
			param: api.CreateUser{
				Username: "eve",
				Password: "456789",
			},
		},
		// user delete
		{
			path:       "/user/delete",
			method:     "POST",
			name:       "not exist",
			expectCode: 1,
			expectErr:  model.UserNotExistErr.Error(),
			param: api.DeleteUser{
				Username: "notexist",
			},
		},
		{
			path:       "/user/delete",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.DeleteUser{
				Username: "eve",
			},
		},
		// role operation
		{
			path:       "/role/create",
			method:     "POST",
			name:       "role name err",
			expectCode: 1,
			expectErr:  model.RoleNameErr.Error(),
			param: api.CreateRole{
				Role: "admin123",
			},
		},
		{
			path:       "/role/create",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.CreateRole{
				Role: "admin",
			},
		},
		{
			path:       "/role/create",
			method:     "POST",
			name:       "role exist",
			expectCode: 1,
			expectErr:  model.RoleExistErr.Error(),
			param: api.CreateRole{
				Role: "admin",
			},
		},
		{
			path:       "/role/create",
			method:     "POST",
			name:       "ok 2",
			expectCode: 0,
			expectErr:  "",
			param: api.CreateRole{
				Role: "ops",
			},
		},
		{
			path:       "/role/create",
			method:     "POST",
			name:       "ok 3",
			expectCode: 0,
			expectErr:  "",
			param: api.CreateRole{
				Role: "dev",
			},
		},
		// role delete
		{
			path:       "/role/delete",
			method:     "POST",
			name:       "not exist",
			expectCode: 1,
			expectErr:  model.RoleNotExistErr.Error(),
			param: api.DeleteRole{
				Role: "notexist",
			},
		},
		{
			path:       "/role/delete",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.DeleteRole{
				Role: "dev",
			},
		},
		// add user role
		{
			path:       "/user/addRole",
			method:     "POST",
			name:       "user not exist",
			expectCode: 1,
			expectErr:  model.UserNotExistErr.Error(),
			param: api.AddUserRole{
				Username: "notexist",
				Role:     "admin",
			},
		},
		{
			path:       "/user/addRole",
			method:     "POST",
			name:       "role not exist",
			expectCode: 1,
			expectErr:  model.RoleNotExistErr.Error(),
			param: api.AddUserRole{
				Username: "bob",
				Role:     "notexist",
			},
		},
		{
			path:       "/user/addRole",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.AddUserRole{
				Username: "bob",
				Role:     "admin",
			},
		},
		{
			path:       "/user/addRole",
			method:     "POST",
			name:       "ok repeat",
			expectCode: 0,
			expectErr:  "",
			param: api.AddUserRole{
				Username: "bob",
				Role:     "admin",
			},
		},
		// auth
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "user not exist",
			expectCode: 1,
			expectErr:  model.UserCheckErr.Error(),
			param: api.Token{
				Username: "notexist",
				Password: "123456",
			},
		},
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "user deleted",
			expectCode: 1,
			expectErr:  model.UserCheckErr.Error(),
			param: api.Token{
				Username: "eve",
				Password: "456789",
			},
		},
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "wrong password",
			expectCode: 1,
			expectErr:  model.UserCheckErr.Error(),
			param: api.Token{
				Username: "bob",
				Password: "random",
			},
		},
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.Token{
				Username: "bob",
				Password: "123456",
			},
		},
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "ok for testing expire",
			expectCode: 0,
			expectErr:  "",
			param: api.Token{
				Username: "bob",
				Password: "123456",
			},
		},
		// check role(with token check)
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "token miss",
			expectCode: 1,
			expectErr:  middleware.TokenRequiredErr.Error(),
			param: api.CheckRole{
				Role: "admin",
			},
		},
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "token invalid",
			expectCode: 1,
			expectErr:  middleware.TokenInvalidErr.Error(),
			param: api.CheckRole{
				Role: "admin",
			},
			header: map[string]*string{"token": invalidToken},
		},
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "ok true",
			expectCode: 0,
			expectErr:  "",
			param: api.CheckRole{
				Role: "admin",
			},
			header: map[string]*string{"token": bobToken},
		},
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "ok false",
			expectCode: 0,
			expectErr:  "",
			param: api.CheckRole{
				Role: "ops",
			},
			header: map[string]*string{"token": bobToken},
		},
		// all roles
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param:      nil,
			header:     map[string]*string{"token": bobToken},
		},
		// Invalidate
		{
			path:       "/auth/logout",
			method:     "POST",
			name:       "invalid token",
			expectCode: 1,
			expectErr:  middleware.TokenInvalidErr.Error(),
			param:      nil,
			header:     map[string]*string{"token": invalidToken},
		},
		{
			path:       "/auth/logout",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param:      nil,
			header:     map[string]*string{"token": bobToken},
		},
		// deleted token
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "invalid token",
			expectCode: 1,
			expectErr:  middleware.TokenInvalidErr.Error(),
			param:      nil,
			header:     map[string]*string{"token": bobToken},
		},
		// expired token
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "expired token",
			expectCode: 1,
			expectErr:  middleware.TokenExpiredErr.Error(),
			param:      nil,
			header:     map[string]*string{"token": bobExpireToken},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.path == "/user/roles" && c.name == "expired token" {
				time.Sleep(time.Duration(model.TokenLifeTime+1) * time.Second) // wait for expire
			}

			w := post(c.path, c.method, c.param, c.header, router)
			var response api.Response
			json.Unmarshal([]byte(w.Body.String()), &response)

			assert.Equal(t, c.expectCode, response.Status)
			assert.Equal(t, c.expectErr, response.Error)

			if c.path == "/auth/token" && c.name == "ok" {
				initToken = w.Header().Get("token")
			}

			if c.path == "/auth/token" && c.name == "ok for testing expire" {
				expireToken = w.Header().Get("token")
			}

			fmt.Printf("Path: [%s] Case: [%s] Param: [%+v] Header [%+v] ExpectCode: [%d] ExpectErr: [%s]\n",
				c.path, c.name, c.param, c.header, c.expectCode, c.expectErr)
		})
	}
}
