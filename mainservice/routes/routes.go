package routes

import (
	"github.com/Amir1848/samrt-library/routes/library"
	"github.com/Amir1848/samrt-library/routes/users"
	"github.com/gin-gonic/gin"
)

func AddMainRoutes(routeEngine *gin.Engine) {
	r := routeEngine.Group("api")

	users.AddRoutes(r)
	library.AddRoutes(r)
}
