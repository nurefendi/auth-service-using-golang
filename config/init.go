package config

import (
	"auth-service/config/database"
	"auth-service/config/environment"
)

func Init() {
	// Get the env profile
	environment.GetEnvironmentConfig()
	// Connect to DB
	database.CreateDBConnection()
}
