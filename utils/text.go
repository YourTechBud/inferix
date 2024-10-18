package utils

import "strings"

// SanitizeJSONText sanitizes the provided JSON text by extracting the JSON object from the text.
func SanitizeJSONText(input string) string {
	// Find the text within the code block if one exists
	if strings.Contains(input, "<code>") {
		// Find the index of the first "<code>"
		start := strings.Index(input, "<code>")
		if start == -1 {
			return ""
		}

		// Find the index of the last "</code>"
		end := strings.LastIndex(input, "</code>")
		if end == -1 {
			return ""
		}

		// Return the substring between the start and end indices
		input = input[start : end+1]
	}

	// Find the index of the first "{"
	start := strings.Index(input, "{")
	if start == -1 {
		return ""
	}

	// Find the index of the last "}"
	end := strings.LastIndex(input, "}")
	if end == -1 {
		return ""
	}

	// Return the substring between the start and end indices
	return input[start : end+1]
}
