package env

import "os"

var defaultValues = map[string]string{
	"GIT_DIR": "./tmp/git",
}

func Get(key string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	val, ok := defaultValues[key]
	if !ok {
		return ""
	}

	return val
}