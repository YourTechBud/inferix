package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DownloadAndConvertToBase64 downloads an image from a URL and converts it to base64
func DownloadAndConvertToBase64(url string) (string, error) {
	// Only handle HTTP(S) URLs
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return url, nil
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(body)

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // Default to JPEG if no content type
	}

	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data), nil
}
