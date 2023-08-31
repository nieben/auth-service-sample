package model

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	TokenLength = 32
)

var (
	TokenLifeTime int64

	Tokens = make(map[string]*Token, 0)

	tLock sync.RWMutex // Tokens lock

	TokenNotExistErr = errors.New("token not exist")
)

type Token struct {
	Token     string `json:"token"`
	User      *User
	CreatedAt int64 `json:"createdAt" binding:"-"`
	ExpireAt  int64 `json:"expireAt"`
}

func GenerateToken(user *User) *Token {
	ts := time.Now().Unix()
	t := &Token{
		User:      user,
		CreatedAt: ts,
		ExpireAt:  ts + TokenLifeTime, // 2h life time
	}
	b, _ := json.Marshal(t.User)
	b = []byte(string(b) + uuid.New().String())
	t.Token = fmt.Sprintf("%x", md5.Sum(b))

	tLock.Lock()
	Tokens[t.Token] = t
	tLock.Unlock()

	return t
}

func GetToken(token string) *Token {
	// no lock
	return Tokens[token]
}

func (t *Token) Remove() error {
	tLock.Lock()
	defer tLock.Unlock()

	if _, ok := Tokens[t.Token]; !ok {
		return TokenNotExistErr
	} else {
		delete(Tokens, t.Token)
	}

	return nil
}
