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
	dev := flag.Bool("dev", false, "development mode")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("failed to load .env: %w", err))
	}
	port := "8000"
	if envport, ok := os.LookupEnv("port"); ok {
		port = envport
	}
	secret := os.Getenv("SECRET_KEY")
	dbUrl := os.Getenv("POSTGRES_DB_URL")

	return Config{
		Prod:   !*dev,
		Port:   port,
		DbURL:  dbUrl,
		Secret: secret,
	}
}
