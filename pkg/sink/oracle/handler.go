package sink

import (
	"context"
	"database/sql/driver"
	"strconv"
	"strings"
	"time"

	goora "github.com/sijms/go-ora/v2"
	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultFlushTime int = 30
)

type OracleSinkHandler struct {
	sink.BaseSink

	sinkOracleCfg config.OracleSinkConfig
	dbConn        *goora.Connection
	ticker        *time.Ticker
}

// SinkName return Oracle sink name
func (o *OracleSinkHandler) SinkName() string { return utils.SinkOracleTagName }

// WriteData write data into Oracle
func (o *OracleSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Oracle-Sink][Current config: " + o.SinkAliasName + "]Start waiting for data to be written...")

	var (
		batch        [][]driver.Value
		batchCounter int
	)

	// flush data function
	flushToDatabase := func(batchInsertData [][]driver.Value) bool {
		sqlText := utils.GenerateInsertSQL(o.sinkOracleCfg.TableName, ":", len(batch))

		for i := 0; i < batchCounter; i++ {
			o.Metrics.OnSinkOutput(o.StreamName, o.SinkAliasName)
		}

		bulkInsertResult, bulkInsertErr := o.dbConn.BulkInsert(sqlText, batchCounter, batchInsertData...)
		if bulkInsertErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Oracle-Sink][Current config: " + o.SinkAliasName + "]batch flush data error, Reason for Exception: " + bulkInsertErr.Error())
			return false
		}
		rowsAffected, _ := bulkInsertResult.RowsAffected()
		logger.Logger.Debug(utils.LogServiceName +
			"[Oracle-Sink][Current config: " + o.SinkAliasName + "]Write " + strconv.Itoa(int(rowsAffected)) + " row data")

		for i := 0; i < int(rowsAffected); i++ {
			o.Metrics.OnSinkOutputSuccess(o.StreamName, o.SinkAliasName)
		}
		batchCounter = 0
		batch = nil

		return true
	}

	for {
		if o.dbConn == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Oracle-Sink][Current config: " + o.SinkAliasName + "]Oracle Connection Closed or Not Connected!!")
			return
		}

		select {
		case <-o.ticker.C:
			logger.Logger.Debug(utils.LogServiceName +
				"[Oracle-Sink][Current config: " + o.SinkAliasName + "]commit on schedule...")
			if batchCounter > 0 && len(batch) > 0 {
				if !flushToDatabase(batch) {
					break
				}
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[Oracle-Sink][Current config: " + o.SinkAliasName + "]batch is empty! batchCounter: " + strconv.Itoa(batchCounter))
			}
		case row, ok := <-o.GetFromTransformChan():
			if !ok || row == nil {
				break
			}
			o.Metrics.OnSinkInput(o.StreamName, o.SinkAliasName)

			if len(batch) == 0 || batchCounter > 0 {
				if len(batch) == 0 {
					batch = make([][]driver.Value, len(row.Data))
				}

				for i, colData := range row.Data {
					if batch[i] == nil {
						batch[i] = make([]driver.Value, 0)
					}
					batch[i] = append(batch[i], colData)
				}

				batchCounter += 1
				o.Metrics.OnSinkInputSuccess(o.StreamName, o.SinkAliasName)
			}

			if len(batch) > 0 && batchCounter == o.sinkOracleCfg.BulkSize {
				// 在一定量以后再提交
				if !flushToDatabase(batch) {
					break
				}
			}
		}
	}
}

// InitSink init Oracle connection
func (o *OracleSinkHandler) InitSink() {
	addressElements := strings.Split(o.sinkOracleCfg.Address, ":")
	connStr := goora.BuildUrl(
		addressElements[0],
		cast.ToInt(addressElements[1]),
		o.sinkOracleCfg.Database,
		o.sinkOracleCfg.Username,
		o.sinkOracleCfg.Password,
		nil)
	conn, connErr := goora.NewConnection(connStr, nil)
	if connErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[Oracle-Sink][Current config: " + o.SinkAliasName + "]Create Oracle connection failed! Reason for Exception: " + connErr.Error())
		return
	}
	openErr := conn.Open()
	if openErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[Oracle-Sink][Current config: " + o.SinkAliasName + "]Connect Oracle failed! Reason for Exception: " + openErr.Error())
		return
	}
	pingErr := conn.Ping(context.Background())
	if pingErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[Oracle-Sink][Current config: " + o.SinkAliasName + "]Failed to ping Oracle! Reason for Exception: " + pingErr.Error())
		return
	}
	o.dbConn = conn
}

// CloseSink close Oracle connection
func (o *OracleSinkHandler) CloseSink() {
	if o.dbConn != nil {
		_ = o.dbConn.Close()
	}
	if o.ticker != nil {
		o.ticker.Stop()
	}
	o.Close()
}

func NewOracleSinkHandler(baseSink sink.BaseSink, sinkOracleCfg config.OracleSinkConfig) *OracleSinkHandler {
	handler := &OracleSinkHandler{
		BaseSink:      baseSink,
		sinkOracleCfg: sinkOracleCfg,
		ticker:        time.NewTicker(time.Duration(defaultFlushTime) * time.Second),
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
