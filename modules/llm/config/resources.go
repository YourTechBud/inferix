package config

import (
	"github.com/YourTechBud/inferix/utils"
)

func GetResourcesInfo() []utils.ResourceInfo {
	return []utils.ResourceInfo{
		{
			Path: "backends",
			New: func() utils.Resource {
				return new(BackendConfig)
			},
		},
		{
			Path: "models",
			New: func() utils.Resource {
				return new(ModelConfig)
			},
		},
	}
}
