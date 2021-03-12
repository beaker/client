package api

// StringPtr returns a pointer to a string.
func StringPtr(value string) *string {
	return &value
}

// BoolPtr returns a pointer to a bool.
func BoolPtr(value bool) *bool {
	return &value
}
