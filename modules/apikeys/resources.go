package apikeys

import "github.com/YourTechBud/inferix/utils"

func GetResourcesInfo() []utils.ResourceInfo {
	return []utils.ResourceInfo{
		{
			Path: "keys",
			Type: utils.ConfigResourceType_Array,
			New: func() utils.Resource {
				return new(APIKey)
			},
			ProtectedFields: []string{"hash", "salt"},
		},
	}
}
