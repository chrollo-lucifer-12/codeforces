package main

import (
	"log"

	"github.com/chrollo-lucifer-12/backend/src/db"
	"github.com/chrollo-lucifer-12/backend/src/routes"
	"github.com/gin-gonic/gin"
)	

func main () {
	router := gin.Default()
	database, err := db.Database()

	if err != nil {
		log.Fatal(err)
	}

	routes.SetUpRoutes(router, database)
	router.Run()
}