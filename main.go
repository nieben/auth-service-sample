package main

import (
	"flag"
	"fmt"
	"github.com/nieben/auth-service-sample/model"
	"github.com/nieben/auth-service-sample/route"
)

var (
	port int64
)

func initFlag() {
	flag.Int64Var(&port, "p", 8080, "serve port")
	flag.Int64Var(&model.TokenLifeTime, "tt", 7200, "token life time(second)")
	flag.Parse()

	if port <= 0 {
		panic(any("invalid port"))
	}
	if model.TokenLifeTime <= 0 {
		panic(any("invalid token lifetime"))
	}
}

// @title auth service sample
// @version 1.0
// @description
// @host 127.0.0.1
// @schemes http
func main() {
	initFlag()
	route.Init().Run(fmt.Sprintf(":%d", port))
}
