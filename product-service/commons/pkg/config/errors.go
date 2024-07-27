package config

import "fmt"

func ErrConfigVariableRequired(name string) error {
	return fmt.Errorf("%v required", name)
}

func ErrConfigVariableEmpty(name string) error {
	return fmt.Errorf("%v environment variable is empty, a value was expected", name)
}
