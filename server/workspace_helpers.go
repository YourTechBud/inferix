package server

import "fmt"

func getWorkspaceKey(tenant, workspace string) string {
	return fmt.Sprintf("%s:::%s", tenant, workspace)
}
