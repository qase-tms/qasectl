package internal

import (
	"strings"
)

func cleanToken(token string) string {
	return strings.TrimSpace(token)
}

func UpdateToken(token string) error {
	return UpdateConfig(func(c Config) Config {
		c.Token = cleanToken(token)

		return c
	})
}
