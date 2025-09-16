package main

import (
	"github.com/dangLuan01/user-manager/internal/app"
	"github.com/dangLuan01/user-manager/internal/config"
)

func main() {
	
	app.LoadEnv()
	
	cfg := config.NewConfig()
	
	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}