package sink

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultFlushTime int = 30
)

type MySQLSinkHandler struct {
	sink.BaseSink

	sinkMySQLCfg config.MySQLSinkConfig
	dbConn       *sql.DB
	ticker       *time.Ticker
}

// SinkName return MySQL sink name
func (m *MySQLSinkHandler) SinkName() string {
	return utils.SinkMySQLTagName
}

// WriteData write data
func (m *MySQLSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Start waiting for data to be written...")

	var (
		batch        [][]any
		batchCounter int
	)

	// flush data function
	flushToDatabase := func() bool {
		// Get the column width, take the first row of data
		sqlText := utils.GenerateInsertSQL(m.sinkMySQLCfg.TableName, "", len(batch[0]))

		for i := 0; i < batchCounter; i++ {
			m.Metrics.OnSinkOutput(m.StreamName, m.SinkAliasName)
		}

		// Open a transaction
		tx, txBeginErr := m.dbConn.Begin()
		if txBeginErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Transaction Startup Exception, Reason for Exception: " + txBeginErr.Error())
			return false
		}

		// Preparing SQL Statements
		stmt, stmtErr := tx.Prepare(sqlText)
		if stmtErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Pre-compile Statement Exception, Reason for Exception: " + stmtErr.Error())
			return false
		}
		defer func() {
			_ = stmt.Close()
		}()

		// Batch insertion of data and counting of affected rows
		totalRowsAffected := int64(0)
		for _, row := range batch {
			execResult, execErr := stmt.Exec(row...)
			if execErr != nil {
				_ = tx.Rollback()
				logger.Logger.Error(utils.LogServiceName +
					"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Statement Execution Exception, Reason for Exception: " + execErr.Error())
				return false
			}

			// Get the number of rows affected by each insertion and add them up
			rowsAffected, err := execResult.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				logger.Logger.Error(utils.LogServiceName +
					"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Get affecting rows exception, reason for exception: " + execErr.Error())
				return false
			}
			totalRowsAffected += rowsAffected
		}

		// Submission of transactions
		if commitErr := tx.Commit(); commitErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[MySQL-Sink][Current config: " + m.SinkAliasName + "]Transaction Commit Exception, Reason for Exception: " + commitErr.Error())
			return false
		}

		for i := 0; i < int(totalRowsAffected); i++ {
			m.Metrics.OnSinkOutputSuccess(m.StreamName, m.SinkAliasName)
		}
		batchCounter = 0
		batch = nil

		return true
	}

	for {
		if m.dbConn == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[MySQL-Sink][Current config: " + m.SinkAliasName + "]MySQL Connection Closed or Not Connected!")
			return
		}

		select {
		case <-m.ticker.C:
			logger.Logger.Debug(utils.LogServiceName +
				"[MySQL-Sink][Current config: " + m.SinkAliasName + "]commit on schedule...")
			if len(batch) > 0 && batchCounter > 0 {
				if !flushToDatabase() {
					break
				}
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[MySQL-Sink][Current config: " + m.SinkAliasName + "]batch is empty! batchCounter: " + strconv.Itoa(batchCounter))
			}
		case row, ok := <-m.GetFromTransformChan():
			if !ok || row == nil {
				break
			}
			m.Metrics.OnSinkInput(m.StreamName, m.SinkAliasName)

			if len(batch) == 0 || batchCounter > 0 {
				if len(batch) > 0 {
					batch = make([][]any, 0)
				}
				batch = append(batch, row.Data)

				batchCounter += 1
				m.Metrics.OnSinkInputSuccess(m.StreamName, m.SinkAliasName)
			}

			if len(batch) > 0 && batchCounter > 0 {
				// 在一定量以后再提交
				if !flushToDatabase() {
					break
				}
			}
		}
	}
}

// InitSink init MySQL connection
func (m *MySQLSinkHandler) InitSink() {
	dbConn, dbOpenErr := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s",
			m.sinkMySQLCfg.Username,
			m.sinkMySQLCfg.Password, // no need to escape,
			m.sinkMySQLCfg.Address,
			m.sinkMySQLCfg.Database))
	if dbOpenErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[MySQL-Sink][Current config: " + m.SinkAliasName + "]Connection to MySQL failed! Reason for error: " + dbOpenErr.Error())
		return
	}
	if pingErr := dbConn.Ping(); pingErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[MySQL-Sink][Current config: " + m.SinkAliasName + "]Ping to MySQL failed! Reason for error: " + pingErr.Error())
		return
	}
	m.dbConn = dbConn
}

// CloseSink close MySQL connection
func (m *MySQLSinkHandler) CloseSink() {
	if m.dbConn != nil {
		_ = m.dbConn.Close()
	}
	if m.ticker != nil {
		m.ticker.Stop()
	}
	m.Close()
}

func NewMySQLSinkHandler(baseSink sink.BaseSink, sinkMySQLCfg config.MySQLSinkConfig) *MySQLSinkHandler {
	handler := &MySQLSinkHandler{
		BaseSink:     baseSink,
		sinkMySQLCfg: sinkMySQLCfg,
		ticker:       time.NewTicker(time.Duration(defaultFlushTime) * time.Second),
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
