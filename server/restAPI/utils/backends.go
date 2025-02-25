package utils

import (
	"errors"
	"io"
	"net/http"
	"strings"

	url "github.com/Dpbm/quantumRestAPI/format"
	logger "github.com/Dpbm/shared/log"
)

func GetBackendsList(pluginName string) (*[]string, error) {
	urlHandler := url.URLForPlugin{PluginName: pluginName}
	response, err := http.Get(urlHandler.GetFullURL())

	if err != nil {
		logger.LogError(err)
		return &[]string{}, err
	}

	if response.StatusCode != 200 {
		responseError := errors.New("failed on get plugin. Status code != 200")
		logger.LogError(responseError)
		return &[]string{}, responseError
	}

	defer response.Body.Close()

	backends, err := io.ReadAll(response.Body)
	if err != nil {
		logger.LogError(err)
		return &[]string{}, err
	}
	lines := strings.Split(string(backends), "\n")

	return &lines, nil
}
