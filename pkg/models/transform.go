package models

type TransformBeforeConvert struct {
	*SourceOutput
	BeforeConvertData any // data before convert
}

type TransformAfterConvert struct {
	*SourceOutput
	AfterConvertData map[string][]any // data after convert
}

type TransformOutput struct {
	*SourceOutput
	Data     []any
	SinkName string
}
