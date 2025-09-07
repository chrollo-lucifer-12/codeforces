package main

import (
	"github.com/chrollo-lucifer-12/backend/src/db"
	"github.com/chrollo-lucifer-12/backend/src/routes"
)

var jwtSecret = []byte("axqiDDn?y|eWEAV")

func main() {
	database := db.InitDB()
	r := routes.SetupRoutes(database)

	r.Run()
}
