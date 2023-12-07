package authutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func AuthorizeOr(requestedRoles []users.UserRole) func(*gin.Context) {
	return func(c *gin.Context) {
		authToken := c.Request.Header["Authorization"]
		if len(authToken) != 1 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("authorization token is not valid"))
		}

		t := strings.Split(authToken[0], " ")
		if len(t) != 2 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("authorization token is not valid"))
		}

		tokenString := t[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			key := os.Getenv("JwtKey")

			return []byte(key), nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			expInUnixDate, ok := claims["exp"].(float64)
			if !ok {
				c.AbortWithError(http.StatusUnauthorized, errors.New("exp date is not valid"))
			}
			exp := time.Unix(int64(expInUnixDate), 0)

			if time.Now().After(exp) {
				c.AbortWithError(http.StatusUnauthorized, err)
			}

			c.Set("userId", int64(claims["userId"].(float64)))

			existedRoles := claims["roles"].([]users.UserRole)
			existedRolesHashset := map[users.UserRole]struct{}{}
			for _, item := range existedRoles {
				existedRolesHashset[item] = struct{}{}
			}

			for _, item := range requestedRoles {
				_, found := existedRolesHashset[item]
				if found {
					return
				}
			}

			c.AbortWithStatus(http.StatusForbidden)
		} else {
			c.AbortWithError(http.StatusUnauthorized, err)
		}

	}
}

func GetUserId(ctx context.Context) (int64, error) {
	userId := ctx.Value("userId").(int64)

	return userId, nil
}
