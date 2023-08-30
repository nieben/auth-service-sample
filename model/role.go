package model

import (
	"errors"
	"sync"
	"time"
)

const (
	RoleRegex = `^[a-zA-Z]{3,15}$`
)

var (
	Roles = make(map[string]*Role, 0)

	rLock sync.RWMutex

	RoleNameErr     = errors.New("role only contains alphabet, len 3-15")
	RoleExistErr    = errors.New("role already exist")
	RoleNotExistErr = errors.New("role not exist")
)

type Role struct {
	Name      string `json:"name" binding:"required"`
	Status    uint8  `json:"status" binding:"-"` // 0: normal 1: deleted
	CreatedAt int64  `json:"createdAt" binding:"-"`
	DeletedAt int64  `json:"deletedAt" binding:"-"`
}

func GetRole(role string) *Role {
	// lock
	rLock.RLock()
	defer rLock.RUnlock()

	return Roles[role]
}

func CreateRole(role string) error {
	rLock.Lock()
	defer rLock.Unlock()

	if _, ok := Roles[role]; ok {
		return RoleExistErr
	} else {
		ts := time.Now().Unix()
		r := &Role{
			Name:      role,
			CreatedAt: ts,
		}
		Roles[r.Name] = r

		return nil
	}
}

func DeleteRole(role string) error {
	rLock.Lock()
	defer rLock.Unlock()

	if _, ok := Roles[role]; !ok {
		return RoleNotExistErr
	} else {
		delete(Roles, role)

		return nil
	}
}
