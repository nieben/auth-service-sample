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

	initBobToken   = ""
	initEveToken   = ""
	expireToken    = ""
	bobToken       = &initBobToken
	eveToken       = &initEveToken
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

// add roles to bob: admin & ops

// now get roles
// bob: [admin ops]
// eve: []

// bob create two token: bob1, bob2
// eve create token: eve1

// bob check role admin => true
// bob check role dev => false

// delete role: ops
// bob check role ops => role not exist

// bob get roles => [admin]

// delete user eve
// eve get roles(with active token eve1) => user not exist
// eve try create token => invalid username or password(deleted)

// invalidate token: bob1
// bob get roles(with bob1) => invalid token

// sleep token lifetime+1 second
// get roles bob(with bob2) => token expired
func TestFullFlow(t *testing.T) {
	cases := []struct {
		path       string
		method     string
		name       string
		expectCode int64
		expectErr  string
		expectData interface{}
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
			name:       "ok bob",
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
			name:       "ok eve",
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
		{
			path:       "/user/addRole",
			method:     "POST",
			name:       "ok ops",
			expectCode: 0,
			expectErr:  "",
			param: api.AddUserRole{
				Username: "bob",
				Role:     "ops",
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
			name:       "ok bob",
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
		{
			path:       "/auth/token",
			method:     "POST",
			name:       "ok eve",
			expectCode: 0,
			expectErr:  "",
			param: api.Token{
				Username: "eve",
				Password: "456789",
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
			expectData: true,
			param: api.CheckRole{
				Role: "admin",
			},
			header: map[string]*string{"token": bobToken},
		},
		// bob does not have dev
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "ok false",
			expectCode: 0,
			expectErr:  "",
			expectData: false,
			param: api.CheckRole{
				Role: "dev",
			},
			header: map[string]*string{"token": bobToken},
		},
		// all roles
		// bob: [admin, ops]
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "ok bob",
			expectCode: 0,
			expectErr:  "",
			expectData: []interface{}{"admin", "ops"},
			param:      nil,
			header:     map[string]*string{"token": bobToken},
		},
		// eve: []
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "ok eve",
			expectCode: 0,
			expectErr:  "",
			expectData: []interface{}{},
			param:      nil,
			header:     map[string]*string{"token": eveToken},
		},
		// delete one added role
		// bob: [admin]  (ops deleted)
		{
			path:       "/role/delete",
			method:     "POST",
			name:       "ok",
			expectCode: 0,
			expectErr:  "",
			param: api.DeleteRole{
				Role: "ops",
			},
		},
		// check deleted role
		{
			path:       "/user/checkRole",
			method:     "POST",
			name:       "check deleted role",
			expectCode: 1,
			expectErr:  model.RoleNotExistErr.Error(),
			param: api.CheckRole{
				Role: "ops",
			},
			header: map[string]*string{"token": bobToken},
		},
		// roles after one role is deleted
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "roles with one been deleted",
			expectCode: 0,
			expectErr:  "",
			expectData: []interface{}{"admin"},
			param:      nil,
			header:     map[string]*string{"token": bobToken},
		},
		// delete user eve
		{
			path:       "/user/delete",
			method:     "POST",
			name:       "ok eve",
			expectCode: 0,
			expectErr:  "",
			param: api.DeleteUser{
				Username: "eve",
			},
		},
		// eveToken become invalid because user eve has been deleted
		// active token with user been deleted
		{
			path:       "/user/roles",
			method:     "POST",
			name:       "active token with user been deleted",
			expectCode: 1,
			expectErr:  model.UserNotExistErr.Error(),
			param:      nil,
			header:     map[string]*string{"token": eveToken},
		},
		// auth with deleted user
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
		// invalid token after bobToken has been invalidated
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
		// wait for token life time second
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
			if c.expectData != nil {
				assert.Equal(t, c.expectData, response.Data)
			}

			if c.path == "/auth/token" && c.name == "ok bob" {
				initBobToken = w.Header().Get("token")
			}

			if c.path == "/auth/token" && c.name == "ok eve" {
				initEveToken = w.Header().Get("token")
			}

			if c.path == "/auth/token" && c.name == "ok for testing expire" {
				expireToken = w.Header().Get("token")
			}

			fmt.Printf("Path: [%s] Case: [%s] Param: [%+v] Header [%+v] ExpectCode: [%d] ExpectErr: [%s]\n",
				c.path, c.name, c.param, c.header, c.expectCode, c.expectErr)
		})
	}
}
