package apikeys

import (
	"github.com/YourTechBud/inferix/utils"
)

type (
	// Config is the configuration for the API keys module
	Config struct {
		APIKeys []APIKey `json:"keys"`
	}

	// APIKey is the struct for an API key
	APIKey struct {
		ID         string `json:"id"`
		Name       string `json:"name" validate:"required,min=3,max=50"`
		Hash       string `json:"hash"`
		Salt       string `json:"salt"`
		LastDigits string `json:"last_digits"`
		Active     bool   `json:"active"` // TODO: Implement this
	}
)

// GetID returns the id of the config
func (a *APIKey) GetID() string {
	return a.ID
}

// Provision generates a new apiKey along with the hash
func (a *APIKey) Provision(ctx *utils.RequestContext) (any, error) {
	// Generate a new id for the key
	a.ID = utils.GenerateID(a.Name)

	// Generate a new api key
	apiKey, hash, err := generateAPIKey(ctx.Workspace(), a.ID)
	if err != nil {
		return nil, err
	}

	// Don't forget to store the hash in the config
	a.Hash = hash
	a.Salt = ""
	a.LastDigits = apiKey[len(apiKey)-4:]

	return map[string]string{"id": a.ID, "key": apiKey}, nil
}

var _ utils.ResourceProvisioner = (*APIKey)(nil)
