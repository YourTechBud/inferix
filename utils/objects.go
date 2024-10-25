package utils

import "strings"

// GetValueAtPath returns the value at the given path in the object
func GetValueAtPath(obj map[string]any, path string, defaultValue any) any {
	// Split the path into parts
	parts := strings.Split(path, "/")

	// Traverse the object
	var next any
	for i, part := range parts {
		// Check if the part exists
		var p bool
		next, p = obj[part]
		if !p {
			return defaultValue
		}

		// Return the value if this is the last part
		if i == len(parts)-1 {
			break
		}

		obj = next.(map[string]any)
	}

	return next
}

// SetValueAtPath sets the value at the given path in the object
func SetValueAtPath(obj map[string]any, path string, value any) {
	// Split the path into parts
	parts := strings.Split(path, "/")

	// Traverse the object
	var next any
	for i, part := range parts {
		// Check if the part exists
		var p bool
		next, p = obj[part]
		if !p {
			// Create a new object
			next = map[string]any{}
			obj[part] = next
		}

		// Create a new object if this is the last part
		if i == len(parts)-1 {
			obj[part] = value
			return
		}

		obj = next.(map[string]any)
	}
}
