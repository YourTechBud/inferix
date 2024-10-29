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

// MergeMaps recursively merges the src map into the dst map
func MergeMaps(dst, src map[string]any) {
	// Return if the source is nil
	if src == nil {
		return
	}

	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		} else {
			// Check if the value is an array
			if srcValue, ok := v.([]any); ok {
				// First check if the destination is an array
				destValue, ok := dst[k].([]any)
				if !ok {
					// Simply set the value
					dst[k] = v
					continue
				}

				// Append the values
				dst[k] = append(destValue, srcValue...)
			} else if srcValue, ok := v.(map[string]any); ok {
				// First check if the destination is a map
				destValue, ok := dst[k].(map[string]any)
				if !ok {
					// Simply set the value
					dst[k] = v
					continue
				}

				MergeMaps(destValue, srcValue)
			} else {
				dst[k] = v
			}
		}
	}
}
