package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/model"
)

func PrintInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if gin.IsDebugging() {
			fmt.Println("Current Users", len(model.Users))
			for _, v := range model.Users {
				fmt.Printf("%+v ", *v)
				fmt.Println()
			}

			fmt.Println("Current Roles", len(model.Roles))
			for _, v := range model.Roles {
				fmt.Printf("%+v ", *v)
				fmt.Println()
			}

			fmt.Println("Current UserRoles", len(model.UserRoles))
			for k, v := range model.UserRoles {
				fmt.Printf("%s: ", k)
				for kk := range v {
					fmt.Printf("%s ", kk)
				}
				fmt.Println()
			}

			fmt.Println("Current Tokens", len(model.Tokens))
			for _, v := range model.Tokens {
				fmt.Printf("%+v\n", *v)
			}
		}
	}
}
