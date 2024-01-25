package library

import (
	"net/http"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	libraryService "github.com/Amir1848/samrt-library/services/library"
	"github.com/Amir1848/samrt-library/services/users"
	"github.com/gin-gonic/gin"
)

type libraryWithItem struct {
	Library      *libraryService.GnrLibrary
	LibraryItems []*libraryService.GnrLibraryItem
}

func AddRoutes(router *gin.RouterGroup) {
	r := router.Group("library")

	r.POST("insert", authutil.AuthorizeOr(users.RoleLibAdmin), func(ctx *gin.Context) {
		lb := libraryWithItem{}
		err := ctx.ShouldBindJSON(&lb)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := libraryService.Insert(ctx, lb.Library, lb.LibraryItems)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	r.GET("/get-library", authutil.AuthorizeOr(users.RoleLibAdmin), func(ctx *gin.Context) {

		result, err := libraryService.GetLibraries(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"libraries": result,
		})
	})
}
