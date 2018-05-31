package tencentcloud

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IntToString converts an integer to a string
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

func StrToInt(value string) int {
	return int(StrToInt64(value))
}

// Int64ToString converts an Int64 integer to a string
func Int64ToStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

func StrToInt64(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

// StrToBool converts a given string specifying "True", or "False" in various lower case or upper case
// into a boolean of true or false.
func StrToBool(value string) bool {
	switch strings.ToUpper(value) {
	case "TRUE":
		return true
	default:
		return false
	}
}

// BoolToStr converts a boolean into a string of either TRUE or FALSE.
func BoolToStr(value bool) string {
	if value {
		return "TRUE"
	}
	return "FALSE"
}

// GetRequiredEnvVars takes a map and reads the values of the keys from the environment.
// If the environment variable doesn't exist, it uses the value from the map itself.
// If the environment variable exists, the map is updated with the value retrieved from the environment.
// This function requires SecretId and SecretKey to be defined in the environment as a special case.
// All other variables are retrieved from the keys defined in the map, then retrieved from the environment using the key as the variable name.
// If the environment variable is empty, it retains the value set originally in the map, which could be empty, or the default value.
// This is used currently only during testing.
func GetRequiredEnvVars(varsRequired map[string]string) {
	varsRequired["SecretId"] = ""
	varsRequired["SecretKey"] = ""
	for k, v := range varsRequired {
		envV := os.Getenv(k)
		if envV == "" && v == "" {
			err := errors.New(fmt.Sprintf("%s not defined in environment", k))
			panic(err)
		}
		if envV != "" {
			varsRequired[k] = envV
		} else {
			varsRequired[k] = v
			// os.Setenv(k, v) // this panics on Windows, so don't call it
		}
	}
}

// GetEnvVar returns a value retrieved from environment for the given name.
func GetEnvVar(name string) string {
	return os.Getenv(name)
}
