package utils

func Coalesce[T any](in *T, def T) T {
	if in == nil {
		return def
	}
	return *in
}
