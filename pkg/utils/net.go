package utils

import "os"

// GetHostName get host name
func GetHostName() string {
	hostName, getHostNameErr := os.Hostname()
	if getHostNameErr != nil {
		return ServiceName
	}
	return hostName
}

// IsDockerRuntime  is in Docker runtime
func IsDockerRuntime() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
