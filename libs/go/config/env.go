package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// loadDotEnv reads a .env file into the OS environment.
// Uses godotenv.Load (not Overload) so pre-existing env vars take precedence.
func loadDotEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("%w: %v", ErrEnvFileLoad, err)
	}
	return nil
}
