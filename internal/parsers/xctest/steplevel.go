package xctest

import "strings"

type stepLevel int

const (
	All = iota
	User
	FirstLevel
)

func parseStepLevel(level string) stepLevel {
	switch strings.ToLower(level) {
	case "all":
		return All
	case "user":
		return User
	default:
		return FirstLevel
	}
}
