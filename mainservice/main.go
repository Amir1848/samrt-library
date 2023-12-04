package main

import (
	"os"

	"github.com/Amir1848/samrt-library/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	var envFileName string
	env := os.Getenv("ENV")
	if env == "development" {
		envFileName = "development.env"
	} else if env == "production" {
		envFileName = "production.env"
	} else {
		//todo: enable panic
		// panic("enviroment variable has not been set")
		envFileName = "development.env"
	}

	err := godotenv.Load(envFileName)
	if err != nil {
		panic("cannot load env file")
	}
}

func main() {
	r := gin.Default()

	routes.AddMainRoutes(r)

	r.Run()
}
