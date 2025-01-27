package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// init() will run only once.
// init() will run after global variable initialization of each package and
// before main() function. init() will only run if the package is imported.
func init() {
	fileName := ".env"
	if fileExists(fileName) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
	}
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetPaymentServiceURL() string {
	return getEnvironmentValue("PAYMENT_SERVICE_URL")
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}

	return os.Getenv(key)
}
