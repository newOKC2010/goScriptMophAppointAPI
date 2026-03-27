package loadenv

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadDBconnec() string {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	return os.Getenv("DB_URL")
}

func LoadScheduleTime() string {
	godotenv.Load()
	return os.Getenv("TIME_SCHEDULE")
}

func LoadDateCount() string {
	godotenv.Load()
	return os.Getenv("DATE_COUNT")
}

func LoadMOPH() (url, clientKey, secretKey string) {
	godotenv.Load()
	return os.Getenv("URL_MOPH"), os.Getenv("CLIENT_VALUE"), os.Getenv("SECRET_VALUE")
}
