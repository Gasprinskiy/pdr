package envloader

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type envloader struct{}

func Init() *envloader {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: ", err)
	}

	return &envloader{}
}

func (*envloader) MustGetString(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("variable %s not found:", key)
	}

	return val
}

func (*envloader) MustGetInt(key string) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("variable %s not found:", key)
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("could not parse %s variable", key)
	}

	return intVal
}
