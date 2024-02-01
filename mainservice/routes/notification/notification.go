package notification

import (
	"net/http"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	"github.com/Amir1848/samrt-library/services/notification"
	"github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.RouterGroup) {
	r := router.Group("notification")

	r.GET("get", authutil.AuthorizeOr(users.RoleLibAdmin, users.RoleLibAdmin, users.RoleUser), func(ctx *gin.Context) {
		var resuult, err = notification.GetUsersNotifications(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resuult)
	})

}
