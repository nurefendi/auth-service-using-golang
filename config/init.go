package config

import (
	"github.com/nurefendi/auth-service-using-golang/config/environment"
)

func Init() {
	// Get the env profile
	environment.GetEnvironmentConfig()
}
