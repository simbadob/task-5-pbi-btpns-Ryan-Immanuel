package main

import (
	"rakamin-final/database"
	"rakamin-final/initializers"
	"rakamin-final/router"
)

func main() {
	initializers.LoadEnv()
	database.Connect()
	database.Migrate()
	r := router.ConfigureRouter()
	r.Run()
}
