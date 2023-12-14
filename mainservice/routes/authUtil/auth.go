package authutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Amir1848/samrt-library/obvious"
	"github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func AuthorizeOr(requestedRoles ...users.UserRole) func(*gin.Context) {
	return func(c *gin.Context) {
		tokenString, err := validateAndGetAuthorizationToken(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		token, err := createJwtToken(tokenString)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			err = validateTokenExpiration(c, claims)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			_, err = setUserId(c, claims)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			if len(requestedRoles) == 0 {
				return
			}

			requestedRolesSet := map[users.UserRole]struct{}{}
			for _, item := range requestedRoles {
				requestedRolesSet[item] = struct{}{}
			}

			existedRoles := claims["roles"].([]interface{})
			for _, existedItem := range existedRoles {
				_, found := requestedRolesSet[users.UserRole(int(existedItem.(float64)))]
				if found {
					return
				}
			}

			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

	}
}

func GetUserId(ctx context.Context) int64 {
	userId := ctx.Value("userId").(int64)

	return userId
}

func validateAndGetAuthorizationToken(c *gin.Context) (string, error) {
	authToken := c.Request.Header["Authorization"]

	if len(authToken) != 1 {
		return "", obvious.ReturnPrefailedCondition(c, "authorization token is not valid")
	}

	t := strings.Split(authToken[0], " ")
	if len(t) != 2 {
		return "", obvious.ReturnPrefailedCondition(c, "authorization token is not valid")
	}

	return t[1], nil
}

func createJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := os.Getenv("JwtKey")

		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func validateTokenExpiration(ctx context.Context, claims jwt.MapClaims) error {
	expInUnixDate, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("expInUnixDate is not valid")
	}
	exp := time.Unix(int64(expInUnixDate), 0)

	if time.Now().After(exp) {
		return obvious.ReturnPrefailedCondition(ctx, "token has expired")
	}

	return nil
}

func setUserId(c *gin.Context, claims jwt.MapClaims) (int64, error) {
	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("user id not found")
	}

	c.Set("userId", int64(userId))
	return int64(userId), nil
}
