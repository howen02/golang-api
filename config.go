package main

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
	DBUser string
	DBPassword string
	DBAddress string
	DBName string
	JWSecret string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		DBUser: getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBAddress: fmt.Sprintf("%s:%s",getEnv("DB_HOST", "127.0.0.1"), 
		getEnv("DB_PORT", "3306")),
		DBName: getEnv("DB_NAME", "alfred"),
		JWSecret: getEnv("JWT_SECRET", "jwtkey"),
		
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {

		return value
	}

	return fallback
}