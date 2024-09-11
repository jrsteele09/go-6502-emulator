package utils

func Value[T any](v *T) T {
	if v == nil {
		v = new(T)
	}
	return *v
}

func Ptr[T any](v T) *T {
	return &v
}
