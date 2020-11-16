package shell

import (
	"fmt"
	"net/url"
	"strings"
)

// Pod是否在线
func CheckQuery(query url.Values) (bool, error) {
	appName 	:= query.Get("appName")
	serverName 	:= query.Get("serverName")
	podName 	:= query.Get("podName")
	if appName == "" || serverName == "" || podName == "" {
		return false, fmt.Errorf("empty params: app: %s, server: %s, pod: %s\n", appName, serverName, podName)
	}

	nodeIP 	:= query.Get("nodeIP")
	history := query.Get("history")
	if strings.ToLower(history) == "true" {
		if nodeIP == "" {
			return false, fmt.Errorf("empty nodeIP can not exec history.\n")
		} else {
			return false, nil
		}
	}

	return true, nil
}
