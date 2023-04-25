package pkg

func Contains[T comparable](values []T, v T) bool {
	for _, value := range values {
		if v == value {
			return true
		}
	}

	return false
}
