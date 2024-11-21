package security

import (
	"github.com/YourTechBud/inferix/utils"
)

type (
	// Config is the configuration for the API keys module
	Config struct {
		APIKeys []APIKey `json:"api-keys"`
	}

	// APIKey is the struct for an API key
	APIKey struct {
		ID         string `json:"id"`
		Desc       string `json:"desc" validate:"required,min=3,max=50"`
		Hash       string `json:"hash"`
		Salt       string `json:"salt"`
		LastDigits string `json:"last_digits"`
		Active     bool   `json:"active"`
	}

	// Metadata is the metadata for the API key resource
	Metadata struct {
		LastUsedAt int64 `json:"last_used_at"`
	}
)

// GetID returns the id of the config
func (a *APIKey) GetID() string {
	return a.ID
}

// Provision generates a new apiKey along with the hash
func (a *APIKey) Provision(ctx *utils.RequestContext) (returningValue, metadata any, err error) {
	// Generate a new api key
	apiKey, hash, err := generateAPIKey(ctx.Tenant(), ctx.Workspace(), a.ID)
	if err != nil {
		return nil, nil, err
	}

	// Don't forget to store the hash in the config
	a.Hash = hash
	a.Salt = ""
	a.LastDigits = apiKey[len(apiKey)-6:]

	return map[string]string{"id": a.ID, "key": apiKey}, Metadata{utils.CurrentTime()}, nil
}

// Update simply updates the title of the previous api key. No other keys is allowed to be updated.
func (a *APIKey) Update(ctx *utils.RequestContext, oldValue, oldMetadata any) (returningValue, metadata any, err error) {
	// Set all the values from the old key
	oldResource := oldValue.(*APIKey)

	a.Hash = oldResource.Hash
	a.Salt = oldResource.Salt
	a.LastDigits = oldResource.LastDigits

	return map[string]string{"id": a.ID}, metadata, nil
}

var _ utils.ResourceProvisioner = (*APIKey)(nil)
var _ utils.ResourceUpdater = (*APIKey)(nil)
