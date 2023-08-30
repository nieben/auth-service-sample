package api

import (
	"github.com/nieben/auth-service-sample/model"
	"regexp"
	"strings"
)

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (in *CreateUser) Check() error {
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	nameReg := regexp.MustCompile(model.UsernameRegex)
	if !nameReg.MatchString(in.Username) {
		return model.UserNameErr
	}
	pwdReg := regexp.MustCompile(model.PwdRegex)
	if !pwdReg.MatchString(in.Password) {
		return model.UserPwdErr
	}

	return nil
}

type DeleteUser struct {
	Username string `json:"username" binding:"required"`
}

func (in *DeleteUser) Check() error {
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	nameReg := regexp.MustCompile(model.UsernameRegex)
	if !nameReg.MatchString(in.Username) {
		return model.UserNameErr
	}

	return nil
}

type CreateRole struct {
	Role string `json:"role" binding:"required"`
}

func (in *CreateRole) Check() error {
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	roleReg := regexp.MustCompile(model.RoleRegex)
	if !roleReg.MatchString(in.Role) {
		return model.RoleNameErr
	}

	return nil
}

type DeleteRole struct {
	Role string `json:"role" binding:"required"`
}

func (in *DeleteRole) Check() error {
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	roleReg := regexp.MustCompile(model.RoleRegex)
	if !roleReg.MatchString(in.Role) {
		return model.RoleNameErr
	}

	return nil
}

type AddUserRole struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (in *AddUserRole) Check() error {
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	nameReg := regexp.MustCompile(model.UsernameRegex)
	if !nameReg.MatchString(in.Username) {
		return model.UserNameErr
	}
	roleReg := regexp.MustCompile(model.RoleRegex)
	if !roleReg.MatchString(in.Role) {
		return model.RoleNameErr
	}

	return nil
}

type CheckRole struct {
	Role string `json:"role" binding:"required"`
}

func (in *CheckRole) Check() error {
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	roleReg := regexp.MustCompile(model.RoleRegex)
	if !roleReg.MatchString(in.Role) {
		return model.RoleNameErr
	}

	return nil
}

type Token struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (in *Token) Check() error {
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	nameReg := regexp.MustCompile(model.UsernameRegex)
	if !nameReg.MatchString(in.Username) {
		return model.UserNameErr
	}
	pwdReg := regexp.MustCompile(model.PwdRegex)
	if !pwdReg.MatchString(in.Password) {
		return model.UserPwdErr
	}

	return nil
}
