package sink

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultFlushTime   int = 30
	defaultMaxOpenConn int = 5
	defaultMaxIdleConn int = 5
)

type ClickhouseSinkHandler struct {
	sink.BaseSink

	ckDriver          driver.Conn
	sinkClickhouseCfg config.ClickhouseSinkConfig
	ticker            *time.Ticker
}

// SinkName return Clickhouse sink name
func (c *ClickhouseSinkHandler) SinkName() string {
	return utils.SinkClickhouseTagName
}

// WriteData write Clickhouse data
func (c *ClickhouseSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]start waiting data writing...")

	var (
		batch        driver.Batch
		err          error
		batchCounter int
	)

	// the function flush to database
	flushToDatabase := func(b driver.Batch) {
		nRows := b.Rows()
		for i := 0; i < nRows; i++ {
			c.Metrics.OnSinkOutput(c.StreamName, c.SinkAliasName)
		}

		sendErr := b.Send()
		if sendErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse send data error! Reason: " + sendErr.Error())
		} else {
			logger.Logger.Debug(utils.LogServiceName +
				"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Writing " + strconv.Itoa(batchCounter) + " rows data")

			for i := 0; i < nRows; i++ {
				c.Metrics.OnSinkOutputSuccess(c.StreamName, c.SinkAliasName)
			}
		}

		batch = nil
		batchCounter = 0
	}

	for {
		if c.ckDriver == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse connection already closed or unconfigured!")
			return
		}

		select {
		case <-c.ticker.C:
			logger.Logger.Debug(utils.LogServiceName + "[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]timing of data submission...")
			if batch != nil {
				flushToDatabase(batch)
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]batch is empty! batchCounter: " + strconv.Itoa(batchCounter))
			}
		case row, ok := <-c.GetFromTransformChan():
			if !ok || row == nil {
				break
			}
			c.Metrics.OnSinkInput(c.StreamName, c.SinkAliasName)

			if batch == nil {
				batch, err = c.ckDriver.PrepareBatch(context.Background(), "INSERT INTO "+c.sinkClickhouseCfg.TableName)
				if err != nil {
					logger.Logger.Error(utils.LogServiceName +
						"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse failed to create Batch! Reason: " + err.Error())
					break
				}
			}

			if batch != nil {
				rowDataAppendErr := batch.Append(row.Data...)
				if rowDataAppendErr != nil {
					logger.Logger.Error(utils.LogServiceName +
						"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse data append error! Reason: " + rowDataAppendErr.Error())
					break
				}

				batchCounter += 1
				c.Metrics.OnSinkInputSuccess(c.StreamName, c.SinkAliasName)
			}

			if batch != nil && c.sinkClickhouseCfg.BulkSize == batchCounter {
				flushToDatabase(batch)
			}
		}
	}
}

// autoCreateTable auto create Clickhouse table
func (c *ClickhouseSinkHandler) autoCreateTable() {
	templateSQL, err := RenderClickhouseCreateTableTemplate("autoCreateTable", c.sinkClickhouseCfg)
	if err != nil {
		return
	}
	if execErr := c.ckDriver.Exec(context.Background(), templateSQL); execErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse auto create table error! Reason: " + execErr.Error())
		return
	}
}

// InitSink initialize Clickhouse connection
func (c *ClickhouseSinkHandler) InitSink() {
	ckDriver, openErr := clickhouse.Open(&clickhouse.Options{
		Addr: strings.Split(c.sinkClickhouseCfg.Address, ","),
		Auth: clickhouse.Auth{
			Database: c.sinkClickhouseCfg.Database,
			Username: c.sinkClickhouseCfg.Username,
			Password: c.sinkClickhouseCfg.Password,
		},
		Debug:            c.DebugMode,
		MaxOpenConns:     defaultMaxOpenConn,
		MaxIdleConns:     defaultMaxIdleConn,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	})
	if openErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse connection error! Reason: " + openErr.Error())
		return
	}

	if err := ckDriver.Ping(context.Background()); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			logger.Logger.Fatal(utils.LogServiceName +
				"[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse connection error! Reason: " + exception.Error())
			return
		}
	}

	c.ckDriver = ckDriver
	if c.sinkClickhouseCfg.IsAutoCreateTable {
		c.autoCreateTable()
	}
	logger.Logger.Info(utils.LogServiceName + "[Clickhouse-Sink][Current config: " + c.SinkAliasName + "]Clickhouse sink initialize successful!")
}

// CloseSink close Clickhouse connection
func (c *ClickhouseSinkHandler) CloseSink() {
	if c.ckDriver != nil {
		_ = c.ckDriver.Close()
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.Close()
}

// NewClickhouseSink create Clickhouse Sink
func NewClickhouseSink(baseSink sink.BaseSink, clickhouseCfg config.ClickhouseSinkConfig) *ClickhouseSinkHandler {
	handler := &ClickhouseSinkHandler{
		BaseSink:          baseSink,
		sinkClickhouseCfg: clickhouseCfg,
		ticker:            time.NewTicker(time.Duration(defaultFlushTime) * time.Second),
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
