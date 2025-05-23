package environment

import (
	"fmt"
	"os"
	"strings"
	"auth-service/common/constants"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/log"
)

func loadEnvVariables(key string) error {
	var err error
	if constants.LOCAL == key {
		err = godotenv.Load(".local.env")
	}
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Application is running on %s environment", key))
	return nil
}
func GetEnvironmentConfig() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	log.Info("Selected env = ", env)
	if err := loadEnvVariables(env); err != nil {
		log.Info(fmt.Sprintf("Failed to load env variables, err=%v", err))
		panic(err)
	}
	log.Info("Application Environment Running profile " + env)
}
