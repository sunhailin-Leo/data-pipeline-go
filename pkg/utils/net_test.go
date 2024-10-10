package utils

import "testing"

func TestGetHostName(t *testing.T) {
	host := GetHostName()
	if host == ServiceName {
		t.Errorf("get host name error!")
	}
}

func TestIsDockerRuntime(t *testing.T) {
	if IsDockerRuntime() {
		t.Errorf("IsDockerRuntime error")
	}
}
