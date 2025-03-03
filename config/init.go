package config

import (
	"auth-service/config/environment"
)

func Init() {
	// Get the env profile
	environment.GetEnvironmentConfig()
}
