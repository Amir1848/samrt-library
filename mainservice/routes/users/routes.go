package users

import (
	"net/http"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	usersService "github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
)

func AddRoutes(routerGroup *gin.RouterGroup) {
	r := routerGroup.Group("user")

	r.POST("register", registerUserFunc(usersService.RoleUser))

	r.GET("get-all-users", authutil.AuthorizeOr(usersService.RoleSysAdmin, usersService.RoleLibAdmin), func(ctx *gin.Context) {
		var result, err = usersService.GetAllUserStudentCodes(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"studentCodes": result})
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

	r.POST("register-library-admin", authutil.AuthorizeOr(usersService.RoleSysAdmin), registerUserFunc(usersService.RoleLibAdmin))

}

func registerUserFunc(role usersService.UserRole) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		c := ctx.Request.Context()

		user := usersService.UserViewModel{}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = usersService.RegisterUser(c, &user, []usersService.UserRole{role})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.AbortWithStatus(http.StatusOK)
	}
}
