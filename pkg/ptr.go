package pkg

func Ptr[T any](v T) *T {
	return &v
}

func Deref[T any](ptr *T) T {
	var v T
	if ptr != nil {
		v = *ptr
	}

	return v
}
