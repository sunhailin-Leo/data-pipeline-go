package sink

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

// mockHttpServerHandler Simulating HTTP server-side handlers
func mockHttpServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		panic(fmt.Sprintf("Expected 'POST' request, got '%s'", r.Method))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read request body: %v", err))
	}
	defer func() {
		_ = r.Body.Close()
	}()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal request body: %v", err))
	}

	expected := map[string]interface{}{"name": "test"}
	if data["name"] != expected["name"] {
		panic(fmt.Sprintf("Expected key 'value', got '%v'", data["key"]))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// mockHttpServer Simulating an HTTP server
func mockHttpServer() (*http.Server, net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", mockHttpServerHandler)
	// Specify the IP and port to listen on
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on 127.0.0.1:8080: %v", err))
	}
	// start HTTP server
	server := &http.Server{Handler: mux}
	return server, listener
}

func TestNewHTTPSinkHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()

	// Mock HTTP Server
	server, listener := mockHttpServer()
	defer func() {
		_ = server.Close()
		_ = listener.Close()
	}()

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// Sink HTTP Test
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "http-1",
		ChanSize:      100,
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkHTTPTagName,
		SinkName: "HTTP-1",
		HTTP: config.HTTPSinkConfig{
			URL:         "http://localhost:8080/test",
			ContentType: utils.HTTPContentTypeJSON,
		},
	}

	httpClient := NewHTTPSinkHandler(base, testSinkConfig.HTTP)
	// Sink Write
	go httpClient.WriteData()
	// channel
	c := httpClient.GetFromTransformChan()
	for i := 1; i < 10; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "http-1",
		}
	}

	// for waiting data insert
	time.Sleep(10 * time.Second)

	httpClient.CloseSink()
}
