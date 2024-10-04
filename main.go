package main

import (
	"context"
	"errors"
	"itspay/internal/app/api"
	"log"
)

func main() {
	app, err := api.NewApp(context.Background())
	if err != nil {
		log.Fatal("create app ", err)
	}

	if err := app.Run(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
