package apikeys

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/YourTechBud/inferix/utils"
	"github.com/YourTechBud/inferix/utils/hash"
)

func generateAPIKey(workspace, id string) (apiKey, hashValue string, err error) {
	// Generate a new api key
	randomString, err := utils.GenerateRandomString(32)
	if err != nil {
		return
	}

	keySection := fmt.Sprintf("%s:::%s:::%s", workspace, id, randomString)
	keySection = base64.StdEncoding.EncodeToString([]byte(keySection))

	apiKey = fmt.Sprintf("ik:%s", keySection)

	// Hash the api key
	hashValue, err = hash.DefaultHasher.Hash(randomString)
	if err != nil {
		return
	}

	return
}

func validate(apiKey string, config *Config) error {
	// Split the apiKey
	apiKeyParts := strings.Split(apiKey, ":")
	if len(apiKeyParts) != 2 {
		return utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
	}

	// Check if the api key is valid
	keySection := apiKeyParts[1]
	keySectionRaw, err := base64.StdEncoding.DecodeString(keySection)
	if err != nil {
		return utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
	}

	keySectionParts := strings.Split(string(keySectionRaw), ":::")
	if len(keySectionParts) != 3 {
		return utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
	}

	id := keySectionParts[1]
	key := keySectionParts[2]
	for _, cfg := range config.APIKeys {
		if cfg.ID != id {
			continue
		}
		if err := hash.DefaultHasher.Compare(cfg.Hash, key); err == nil {
			// The key is valid
			return nil
		}
	}

	// The key is invalid
	return utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
}
