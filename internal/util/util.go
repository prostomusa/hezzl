package util

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	return val
}

func GetIntEnv(key string, defaultValue int) (int, error) {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue, nil
	}
	intVal, err := strconv.Atoi(val)
	return intVal, err
}
