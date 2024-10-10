package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/stream"
)

func main() {
	// initialize stream handler
	streamSrv := stream.NewStreamHandler()
	// start service
	streamSrv.Start()
}
