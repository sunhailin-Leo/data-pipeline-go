package sink

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type LocalFileSinkHandler struct {
	sink.BaseSink

	sinkLocalFileCfg config.LocalFileSinkConfig
	file             *os.File
	csvWriter        *csv.Writer
	textWriter       *bufio.Writer
	parquetFp        source.ParquetFile
	parquetWriter    *writer.ParquetWriter
}

func (l *LocalFileSinkHandler) SinkName() string {
	return utils.SinkLocalFileTagName
}

func (l *LocalFileSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Start sending data...")
	for {
		switch l.sinkLocalFileCfg.FileFormatType {
		case utils.LocalFileCSVFormatType, utils.LocalFileTextFormatType:
			if l.file == nil {
				logger.Logger.Error(utils.LogServiceName +
					"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]File closed or not open!")
				return
			}
		case utils.LocalFileParquetFormatType:
			if l.parquetFp == nil {
				logger.Logger.Error(utils.LogServiceName +
					"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]parquet file is closed or not open!")
				return
			}
		default:
			logger.Logger.Error(utils.LogServiceName + "[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Unknown file format!")
			return
		}

		data, ok := <-l.GetFromTransformChan()
		l.Metrics.OnSinkInput(l.StreamName, l.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName +
				"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Data source closed!")
			return
		}
		l.Metrics.OnSinkInputSuccess(l.StreamName, l.SinkAliasName)
		l.Metrics.OnSinkOutput(l.StreamName, l.SinkAliasName)

		switch l.sinkLocalFileCfg.FileFormatType {
		case utils.LocalFileCSVFormatType:
			writeErr := l.csvWriter.Write(utils.InterfaceSliceToStringSlice(data.Data))
			if writeErr != nil {
				logger.Logger.Error(utils.LogServiceName +
					"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]csv file write failed! Reason for error: " + writeErr.Error())
				return
			}
			l.csvWriter.Flush()
			l.Metrics.OnSinkOutputSuccess(l.StreamName, l.SinkAliasName)
		case utils.LocalFileTextFormatType:
			for _, row := range utils.InterfaceSliceToStringSlice(data.Data) {
				_, writeErr := l.textWriter.WriteString(row + l.sinkLocalFileCfg.RowDelimiter)
				if writeErr != nil {
					logger.Logger.Error(utils.LogServiceName +
						"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]text File write failed! Reason for error: " + writeErr.Error())
					return
				}
				flushErr := l.textWriter.Flush()
				if flushErr != nil {
					logger.Logger.Error(utils.LogServiceName +
						"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]text File flush failed! Reason for error: " + flushErr.Error())
					return
				}
			}
			l.Metrics.OnSinkOutputSuccess(l.StreamName, l.SinkAliasName)
		case utils.LocalFileParquetFormatType:
			// TODO 暂时无法写入
			for _, row := range data.Data {
				if writeErr := l.parquetWriter.Write(row); writeErr != nil {
					logger.Logger.Error(utils.LogServiceName +
						"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]parquet file write failed! Reason for error: " + writeErr.Error())
					return
				}
			}
			l.Metrics.OnSinkOutputSuccess(l.StreamName, l.SinkAliasName)
		default:
			logger.Logger.Error(utils.LogServiceName + "[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Unknown file format!")
			return
		}
	}
}

// InitSink initializes the local file sink
func (l *LocalFileSinkHandler) InitSink() {
	switch l.sinkLocalFileCfg.FileFormatType {
	case utils.LocalFileCSVFormatType:
		file, err := os.OpenFile(l.sinkLocalFileCfg.FileName+".csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Failed to open csv file! Reason for error: " + err.Error())
			return
		}
		csvWriter := csv.NewWriter(file)
		l.file = file
		l.csvWriter = csvWriter
	case utils.LocalFileTextFormatType:
		file, err := os.OpenFile(l.sinkLocalFileCfg.FileName+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Failed to open txt file! Reason for error: " + err.Error())
			return
		}
		l.file = file
		l.textWriter = bufio.NewWriter(file)
	case utils.LocalFileParquetFormatType:
		fw, newFwErr := local.NewLocalFileWriter(l.sinkLocalFileCfg.FileName + ".parquet")
		if newFwErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Failed to create parquet file! Reason for error: " + newFwErr.Error())
			return
		}

		pw, newPwErr := writer.NewParquetWriter(fw, nil, 1)
		if newPwErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Failed to create parquet writer! Reason for error: " + newPwErr.Error())
			return
		}

		l.parquetFp = fw
		l.parquetWriter = pw
	default:
		logger.Logger.Error(utils.LogServiceName + "[LocalFile-Sink][Current config: " + l.SinkAliasName + "]Unknown file format!")
		return
	}
}

// CloseSink closes the local file sink
func (l *LocalFileSinkHandler) CloseSink() {
	if l.file != nil {
		_ = l.file.Close()
	}
	if l.parquetFp != nil {
		_ = l.parquetFp.Close()
	}
	if l.parquetWriter != nil {
		_ = l.parquetWriter.WriteStop()
	}
	l.Close()
}

// NewLocalFileHandler creates a new local file sink handler
func NewLocalFileHandler(baseSink sink.BaseSink, sinkLocalFileCfg config.LocalFileSinkConfig) *LocalFileSinkHandler {
	handler := &LocalFileSinkHandler{BaseSink: baseSink, sinkLocalFileCfg: sinkLocalFileCfg}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
