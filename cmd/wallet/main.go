package main

import (
	"log"

	"wallet/config"
	"wallet/internal/app"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	envFilename = "config-local.env"
)

//	@title			wallet api application
//	@version		0.0.1
//	@description	This is a sample wallet api service

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	cfg := new(config.Config)

	if err := godotenv.Load(envFilename); err != nil {
		log.Println("No .env file found")
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
