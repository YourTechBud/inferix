package apikeys

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/YourTechBud/inferix/utils"
	"github.com/YourTechBud/inferix/utils/hash"
)

func GetKeySegments(apiKey string) (tenant, workspace, id, key string, err error) {
	// Split the apiKey
	apiKeyParts := strings.Split(apiKey, ":")
	if len(apiKeyParts) != 2 {
		err = utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
		return
	}

	// Check if the api key is valid
	keySection := apiKeyParts[1]
	keySectionRaw, err := base64.StdEncoding.DecodeString(keySection)
	if err != nil {
		err = utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
		return
	}

	keySectionParts := strings.Split(string(keySectionRaw), ":::")
	if len(keySectionParts) != 4 {
		err = utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
		return
	}

	tenant = keySectionParts[0]
	workspace = keySectionParts[1]
	id = keySectionParts[2]
	key = keySectionParts[3]

	return
}

func generateAPIKey(tenant, workspace, id string) (apiKey, hashValue string, err error) {
	// Generate a new api key
	randomString, err := utils.GenerateRandomString(32)
	if err != nil {
		return
	}

	keySection := fmt.Sprintf("%s:::%s:::%s:::%s", tenant, workspace, id, randomString)
	keySection = base64.StdEncoding.EncodeToString([]byte(keySection))

	apiKey = fmt.Sprintf("ik:%s", keySection)

	// Hash the api key
	hashValue, err = hash.DefaultHasher.Hash(randomString)
	if err != nil {
		return
	}

	return
}

func (module *Module) validate(apiKey string) (string, error) {
	// Get the key segments
	_, _, id, key, err := GetKeySegments(apiKey)
	if err != nil {
		return "", err
	}

	// Check if key exists
	cfg, p := module.apiKeys[id]
	if !p || !cfg.Active {
		return "", utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
	}

	if err := hash.DefaultHasher.Compare(cfg.Hash, key); err != nil {
		// The key is valid
		return "", utils.NewStandardError(http.StatusUnauthorized, "Invalid API key", "invalid_api_key")
	}

	return id, nil
}
