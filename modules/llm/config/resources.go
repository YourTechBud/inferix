package config

import (
	"github.com/YourTechBud/inferix/utils"
)

func GetResourcesInfo() []utils.ResourceInfo {
	return []utils.ResourceInfo{
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
