package env

import "os"

type EnvVars struct {
	Dev bool
}

const (
	devFlag = "DEV"
)

func InitEnvVars() *EnvVars {
	envVars := &EnvVars{}
	isDev, ok := os.LookupEnv(devFlag)
	if ok {
		envVars.Dev = isDev == "1"
	}
	return envVars
}

func (e *EnvVars) IsDev() bool {
	return e.Dev
}

// Return either http or https depending on whether environment is dev or prod
func (e *EnvVars) GetHttpProtocol() string {
	if e.Dev {
		return "http"
	}
	return "https"
}

// Check if dev environment directly from env var
func IsDevDirect() bool {
	isDev, ok := os.LookupEnv(devFlag)
	if ok {
		return isDev == "1"
	}
	return false
}
