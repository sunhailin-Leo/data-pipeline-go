package sink

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultFlushTime int = 30
)

type PostgresSQLHandler struct {
	sink.BaseSink

	sinkPostgresSQLCfg config.PostgresSQLSinkConfig
	pgConn             *pgx.Conn
	ticker             *time.Ticker
}

// SinkName return PostgresSQL sink name
func (p *PostgresSQLHandler) SinkName() string {
	return utils.SinkPostgresSQLTagName
}

// WriteData write data into PostgresSQL
func (p *PostgresSQLHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Start waiting for data to be written...")

	var (
		batch        *pgx.Batch
		batchCounter int
	)

	// flush data function
	flushToDatabase := func(b *pgx.Batch) {
		nRows := b.Len()
		for i := 0; i < nRows; i++ {
			p.Metrics.OnSinkOutput(p.StreamName, p.SinkAliasName)
		}

		// start transaction
		tx, txErr := p.pgConn.Begin(context.Background())
		if txErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Failed to start PostgresSQL transaction! Reason for Exception: " + txErr.Error())
			return
		}

		// create batch
		br := tx.SendBatch(context.Background(), batch)

		// execution batch
		_, ctErr := br.Exec()
		if ctErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]PostgresSQL failed to perform a bulk insert! Reason for Exception: " + ctErr.Error())
			_ = tx.Rollback(context.Background())
			return
		}

		// close batch
		_ = br.Close()

		// commit transaction
		if commitErr := tx.Commit(context.Background()); commitErr == nil {
			for i := 0; i < nRows; i++ {
				p.Metrics.OnSinkOutputSuccess(p.StreamName, p.SinkAliasName)
			}
			logger.Logger.Debug(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Write " + strconv.Itoa(batchCounter) + " row data")
		} else {
			logger.Logger.Error(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Failed to commit transaction! Reason for Exception: " + commitErr.Error())
		}

		batch = nil
		batchCounter = 0
	}

	for {
		if p.pgConn == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]PostgresSQL connection closed or not connected!")
			return
		}

		select {
		case <-p.ticker.C:
			logger.Logger.Debug(utils.LogServiceName +
				"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]commit on schedule...")
			if batch != nil {
				flushToDatabase(batch)
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]batch is empty! batchCounter: " + strconv.Itoa(batchCounter))
			}
		case row, ok := <-p.GetFromTransformChan():
			if !ok || row == nil {
				break
			}
			p.Metrics.OnSinkInput(p.StreamName, p.SinkAliasName)

			if batch == nil {
				batch = &pgx.Batch{}
			}

			if batch != nil {
				symbols := make([]string, len(row.Data))
				for i := range row.Data {
					symbols[i] = "$" + strconv.Itoa(i+1)
				}
				batch.Queue("INSERT INTO "+p.sinkPostgresSQLCfg.TableName+" VALUES ("+strings.Join(symbols, ", ")+")", row.Data...)

				batchCounter += 1
				p.Metrics.OnSinkInputSuccess(p.StreamName, p.SinkAliasName)
			}

			if batch != nil && batchCounter == p.sinkPostgresSQLCfg.BulkSize {
				flushToDatabase(batch)
			}
		}
	}
}

// InitSink init PostgresSQL connection
func (p *PostgresSQLHandler) InitSink() {
	conn, connErr := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s/%s",
			p.sinkPostgresSQLCfg.Username,
			p.sinkPostgresSQLCfg.Password,
			p.sinkPostgresSQLCfg.Address,
			p.sinkPostgresSQLCfg.Database))
	if connErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Failed to connect PostgresSQL! Reason for exception: " + connErr.Error())
		return
	}

	if err := conn.Ping(context.Background()); err != nil {
		logger.Logger.Fatal(utils.LogServiceName + "[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Failed to ping PostgresSQL! Reason for exception: " + err.Error())
		return
	}

	p.pgConn = conn
	logger.Logger.Info(utils.LogServiceName + "[PostgresSQL-Sink][Current config: " + p.SinkAliasName + "]Connect PostgresSQL Successful!")
}

// CloseSink close PostgresSQL connection
func (p *PostgresSQLHandler) CloseSink() {
	if p.pgConn != nil {
		_ = p.pgConn.Close(context.Background())
	}
	if p.ticker != nil {
		p.ticker.Stop()
	}
	p.Close()
}

func NewPostgresSQLHandler(baseSink sink.BaseSink, sinkPostgresSQLCfg config.PostgresSQLSinkConfig) *PostgresSQLHandler {
	handler := &PostgresSQLHandler{
		BaseSink:           baseSink,
		sinkPostgresSQLCfg: sinkPostgresSQLCfg,
		ticker:             time.NewTicker(time.Duration(defaultFlushTime) * time.Second),
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
