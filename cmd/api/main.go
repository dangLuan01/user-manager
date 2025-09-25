package main

import (
	"log"

	"github.com/dangLuan01/user-manager/internal/app"
	"github.com/dangLuan01/user-manager/internal/config"
)

func main() {
	
	app.LoadEnv()
	
	cfg := config.NewConfig()
	
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("Error start app:%s", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Error run app:%s", err)
	}
}