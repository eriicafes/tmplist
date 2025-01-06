package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Prod bool
	Port string
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getConfig() Config {
	prod := flag.Bool("prod", false, "production mode")
	port := flag.String("port", "8000", "application port")
	flag.Parse()

	return Config{
		Prod: *prod,
		Port: *port,
	}
}
