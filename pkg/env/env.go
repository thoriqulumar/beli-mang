package env

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
)

// Env list
const (
	DevelopmentEnv = "development"
	ProductionEnv  = "production"
)

// Env related var
var (
	envName   = "ENV" // can rename this if you want customENV
	goVersion string
)

func init() {
	// env package will read .env file when application is started
	err := SetFromEnvFile(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("failed to set env file: %v\n", err)
	}
	goVersion = runtime.Version()
}

// SetFromEnvFile read env file and set the environment variables
func SetFromEnvFile(filepath string) error {
	if _, err := os.Stat(filepath); err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return err
	}
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		vars := strings.SplitN(text, "=", 2)
		if len(vars) < 2 {
			return err
		}
		if err := os.Setenv(vars[0], vars[1]); err != nil {
			return err
		}
	}
	return nil
}

// ServiceEnv return ENV service environment
func ServiceEnv() string {
	e := os.Getenv(envName)
	if e != "" {
		return e
	}
	return DevelopmentEnv
}

// GoVersion to return current build go version
func GoVersion() string {
	return goVersion
}

// IsDevelopment return true when env is "development"
func IsDevelopment() bool {
	return ServiceEnv() == DevelopmentEnv
}

// IsProduction return true when env is "production"
func IsProduction() bool {
	return ServiceEnv() == ProductionEnv
}
