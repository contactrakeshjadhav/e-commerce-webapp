package config

import "os"

const (
	ecommerceDomain = "ecommerce_DOMAIN"
)

func GetecommerceDomain() (string, error) {
	ecommerceDomain, err := loadEnvironmentVariable(ecommerceDomain)
	if err != nil {
		return "", err
	}
	return ecommerceDomain, nil
}

// loadEnvironmentVariable returns non-empty value when given variable is set and non-empty.
// It returns error on missing variable or empty variable.
func loadEnvironmentVariable(env string) (string, error) {
	value, found := os.LookupEnv(env)
	if !found {
		return "", ErrConfigVariableRequired(env)
	}
	if value == "" {
		return "", ErrConfigVariableEmpty(env)
	}
	return value, nil
}
