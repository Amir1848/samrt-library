package library

import (
	"net/http"

	libraryService "github.com/Amir1848/samrt-library/services/library"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.RouterGroup) {
	r := router.Group("library")

	r.POST("insert", func(ctx *gin.Context) {
		lb := libraryService.GnrLibrary{}
		err := ctx.ShouldBindJSON(&lb)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c := ctx.Request.Context()
		id, err := libraryService.Insert(c, &lb)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	r.GET("/get-libraries", func(ctx *gin.Context) {
		c := ctx.Request.Context()
		result, err := libraryService.GetLibraries(c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"libraries": result,
		})
	})
}
