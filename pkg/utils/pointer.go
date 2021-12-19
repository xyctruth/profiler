package utils

// Bool returns a pointer to a bool.
func Bool(b bool) *bool {
	return &b
}

var BoolPtr = Bool
