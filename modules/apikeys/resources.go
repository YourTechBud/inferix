package apikeys

import "github.com/YourTechBud/inferix/utils"

func GetResourcesInfo() []utils.ResourceInfo {
	return []utils.ResourceInfo{
		{
			Path: "keys",
			New: func() utils.Resource {
				return new(APIKey)
			},
			ProtectedFields: []string{"hash", "salt"},
		},
	}
}
