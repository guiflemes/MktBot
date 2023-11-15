package settings

import (
	"log"
	"os"
)

func GETENV(key string) string {
	env, isPresent := os.LookupEnv(key)

	if !isPresent {
		log.Fatalf("%s is not present in environment variables", key)
	}
	return env
}
