package server

import (
	"fmt"
	"path"

	"github.com/valyala/fastjson"

	"github.com/YourTechBud/inferix/modules/llm"
	llmconfig "github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/security"
	"github.com/YourTechBud/inferix/utils"
)

func getWorkspaceKey(tenant, workspace string) string {
	return fmt.Sprintf("%s:::%s", tenant, workspace)
}

func removeFields(value *fastjson.Value, fields []string) {
	if value == nil {
		return
	}

	switch value.Type() {
	case fastjson.TypeObject:
		for _, field := range fields {
			value.Del(field)
		}
	case fastjson.TypeArray:
		arr := value.GetArray()
		for _, v := range arr {
			removeFields(v, fields)
		}
	}
}

func getConfigDir(configPath, tenant, workspace string) string {
	return path.Join(configPath, tenant, workspace)
}

func getConfigPath(configPath, tenant, workspace string) string {
	return path.Join(getConfigDir(configPath, tenant, workspace), "config.db")
}

var modulesList = []utils.ModuleInfo{
	{
		Name:             "security",
		New:              security.New,
		GetResourcesInfo: security.GetResourcesInfo,
	},
	{
		Name:             "llm",
		New:              llm.New,
		GetResourcesInfo: llmconfig.GetResourcesInfo, // TODO: This should be in the llm module.
	},
}
