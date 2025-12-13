package utils

func Coalesce[T comparable](in *T, def T) T {
	if in == nil {
		return def
	}
	return *in
}
