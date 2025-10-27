package util

func ToPtr[T any](value T) *T {
	return &value
}
