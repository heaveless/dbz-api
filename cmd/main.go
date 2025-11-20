package main

import (
	"log"

	"github.com/heaveless/dbz-api/internal/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	defer app.CloseDbConnection()

	err := app.Svr.Run(":" + env.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
