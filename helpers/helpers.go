package helpers

import (
	"fmt"
	"os"
)

func GetEnvDefault(env, def string) string {
	val := os.Getenv(env)
	if val == "" {
		val = def
	}
	return val
}

func HandleError(e error, f string, fatal bool) {
	if e != nil {
		fmt.Printf("Error: %s during %s", e, f)
	}
	if fatal {
		os.Exit(1)
	}
}
