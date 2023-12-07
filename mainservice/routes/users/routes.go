package users

import (
	"net/http"

	usersService "github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
)

func AddRoutes(routerGroup *gin.RouterGroup) {
	r := routerGroup.Group("user")

	r.POST("register", func(ctx *gin.Context) {
		c := ctx.Request.Context()

		user := usersService.UserViewModel{}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = usersService.RegisterUser(c, &user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.AbortWithStatus(http.StatusOK)
	})

	r.POST("login", func(ctx *gin.Context) {
		c := ctx.Request.Context()

		user := usersService.UserViewModel{}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, res, err := usersService.Login(c, &user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if res {
			ctx.JSON(http.StatusOK, gin.H{"result": token})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"result": res})
		}
	})

}
