package utils

import "testing"

func TestGetHostName(t *testing.T) {
	host := GetHostName()
	if host == ServiceName {
		t.Errorf("get host name error!")
	}
	if host == "" {
		t.Errorf("get host name should not be empty")
	}
}

func TestIsDockerRuntime(t *testing.T) {
	if IsDockerRuntime() {
		t.Errorf("IsDockerRuntime error")
	}
}

func TestGetHostName_Success(t *testing.T) {
	host := GetHostName()
	if host == "" {
		t.Errorf("GetHostName() returned empty string")
	}
	if host == ServiceName {
		t.Errorf("GetHostName() returned ServiceName, indicating an error occurred")
	}
}

func TestIsDockerRuntime_NotDocker(t *testing.T) {
	if IsDockerRuntime() {
		t.Errorf("IsDockerRuntime() returned true in non-Docker environment")
	}
}
