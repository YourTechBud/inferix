package config

import (
	"github.com/YourTechBud/inferix/utils"
)

func GetConfigurationResources() []utils.ResourceConfiguration {
	return []utils.ResourceConfiguration{
		{
			Path: "backends",
			Type: utils.ConfigResourceType_Array,
			New: func() utils.Resource {
				return new(BackendConfig)
			},
		},
		{
			Path: "models",
			Type: utils.ConfigResourceType_Array,
			New: func() utils.Resource {
				return new(ModelConfig)
			},
		},
	}
}
