package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Prod   bool
	Port   string
	DbURL  string
	Secret string
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getConfig() Config {
	prod := flag.Bool("prod", false, "production mode")
	port := flag.String("port", "8000", "application port")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("failed to load .env: %w", err))
	}
	dbUrl := os.Getenv("POSTGRES_DB_URL")
	secret := os.Getenv("SECRET_KEY")

	return Config{
		Prod:   *prod,
		Port:   *port,
		DbURL:  dbUrl,
		Secret: secret,
	}
}
