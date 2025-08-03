package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return intVal
}
