package utils

// NonNilSlice returns an empty slice instead of nil so JSON/Swagger serialize as [].
func NonNilSlice[T any](v []T) []T {
	if v == nil {
		return []T{}
	}
	return v
}
