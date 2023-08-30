package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

const (
	UsernameRegex = `^[a-zA-Z0-9]{3,15}$`
	PwdRegex      = `^[ -~]{6,20}$`
)

var (
	Users     = make(map[string]*User, 0)
	UserRoles = make(map[string]map[string]struct{}, 0)

	uLock  sync.RWMutex // Users lock
	urLock sync.RWMutex // UserRoles lock

	UserNameErr     = errors.New("username only contains number and alphabet, len 3-15")
	UserPwdErr      = errors.New("password only contains ascii space-~, len 6-20")
	UserCheckErr    = errors.New("invalid username or password")
	UserExistErr    = errors.New("user already exist")
	UserNotExistErr = errors.New("user not exist")
)

type User struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Status    uint8  `json:"status" binding:"-"` // 0: normal 1: deleted
	CreatedAt int64  `json:"createdAt" binding:"-"`
	DeletedAt int64  `json:"deletedAt" binding:"-"`
}

func (u *User) CheckPwd(password string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetUser(username string) *User {
	// lock
	uLock.RLock()
	defer uLock.RUnlock()

	return Users[username]
}

func CreateUser(username, password string) error {
	// lock
	uLock.Lock()
	defer uLock.Unlock()

	if _, ok := Users[username]; ok {
		return UserExistErr
	} else {
		ts := time.Now().Unix()
		u := &User{
			Username:  username,
			CreatedAt: ts,
		}

		// encrypt password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)

		Users[u.Username] = u

		return nil
	}
}

func DeleteUser(username string) error {
	// lock
	uLock.Lock()
	defer uLock.Unlock()

	if _, ok := Users[username]; !ok {
		return UserNotExistErr
	} else {
		delete(Users, username)
	}

	return nil
}

func (u *User) AddRole(role *Role) error {
	urLock.Lock()
	defer urLock.Unlock()

	if _, ok := UserRoles[u.Username]; !ok {
		roles := make(map[string]struct{}, 0)
		roles[role.Name] = struct{}{}
		UserRoles[u.Username] = roles
	} else {
		UserRoles[u.Username][role.Name] = struct{}{}
	}

	return nil
}

func (u *User) CheckRole(role string) bool {
	urLock.RLock()
	defer urLock.RUnlock()

	if m, ok := UserRoles[u.Username]; !ok {
		return false
	} else {
		if _, ok := m[role]; ok {
			return true
		} else {
			return false
		}
	}
}

func (u *User) Roles() []string {
	urLock.RLock()
	defer urLock.RUnlock()

	roles := make([]string, 0)
	if m, ok := UserRoles[u.Username]; !ok {
		return roles
	} else {
		for r, _ := range m {
			roles = append(roles, r)
		}
	}

	return roles
}
