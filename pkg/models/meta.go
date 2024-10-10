package models

type MetaData struct {
	StreamName    string // stream name
	SourceTagName string // Source Tag name
	AliasName     string // Source alias name
	SourceObj     any    // Source obj like Kafka, Pulsar .etc
}
