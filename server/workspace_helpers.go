package server

import (
	"fmt"

	"github.com/valyala/fastjson"

	"github.com/YourTechBud/inferix/modules/llm"
	llmconfig "github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/utils"
)

func getWorkspaceKey(tenant, workspace string) string {
	return fmt.Sprintf("%s:::%s", tenant, workspace)
}

func removeFields(value *fastjson.Value, fields []string) {
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

var modulesMap = map[string]utils.ModuleInfo{
	"llm": {
		New:              llm.New,
		GetResourcesInfo: llmconfig.GetResourcesInfo,
	},
}
