package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetString reads an environment variable or returns defaultVal if not set or empty
func GetString(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	return val
}

// GetInt reads an integer environment variable or returns defaultVal
func GetInt(key string, defaultVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return defaultVal
	}
	return parsed
}

// GetInt32 reads an int32 environment variable or returns defaultVal
func GetInt32(key string, defaultVal int32) int32 {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 32)
	if err != nil {
		return defaultVal
	}
	return int32(parsed)
}

// GetInt64 reads an int64 environment variable or returns defaultVal
func GetInt64(key string, defaultVal int64) int64 {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// GetBool reads a boolean environment variable (1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)
func GetBool(key string, defaultVal bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(val))
	if err != nil {
		return defaultVal
	}
	return parsed
}

// GetDuration reads a time.Duration environment variable (e.g. "10s", "1h", "30m")
func GetDuration(key string, defaultVal time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parsed, err := time.ParseDuration(strings.TrimSpace(val))
	if err != nil {
		return defaultVal
	}
	return parsed
}

// GetStringSlice reads a comma-separated list of strings from an environment variable
func GetStringSlice(key string, defaultVal []string) []string {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return defaultVal
	}
	parts := strings.Split(val, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	if len(res) == 0 {
		return defaultVal
	}
	return res
}

// GetRequired returns the string value of an environment variable, or an error if missing/empty
func GetRequired(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return val, nil
}

// LoadEnvFile loads key-value pairs from .env formatted files into the process environment
// It will not overwrite variables that are already set in the environment.
func LoadEnvFile(filenames ...string) error {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Ignore empty lines and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Split by first '='
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Strip quotes if present
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}

			// Strip inline comments if not quoted
			if idx := strings.Index(val, " #"); idx != -1 {
				val = strings.TrimSpace(val[:idx])
			}

			// Set only if not already defined
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}

		_ = file.Close()
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
