package format

import (
	"fmt"
	"strings"

	constants "github.com/Dpbm/quantumRestAPI/constants"
)

type URLForPlugin struct {
	PluginName string
}

func mapNameToPipFormat(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

func (url *URLForPlugin) GetFullURL() string {
	pluginName := url.PluginName
	nameMappedToPIP := mapNameToPipFormat(pluginName)

	return fmt.Sprintf("%s/%s/refs/heads/main/%s/backends.txt", constants.REPO_URL, pluginName, nameMappedToPIP)
}
