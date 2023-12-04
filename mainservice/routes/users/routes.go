package users

import (
	usersService "github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.RouterGroup) {

	r.GET("get-users", func(ctx *gin.Context) {
		usersService.GetUsers(ctx)

		ctx.JSON(200, gin.H{
			"amir": "sahand",
		})
	})

	// r.POST("insert-user", func(ctx *gin.Context) {

	// })
}
