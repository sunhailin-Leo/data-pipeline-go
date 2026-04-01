package sink

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewPostgresSQLHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	pgAddr := testutil.GetEnvOrDefault(testutil.EnvPostgresAddr, "localhost:5432")
	pgUser := testutil.GetEnvOrDefault(testutil.EnvPostgresUser, "testuser")
	pgPass := testutil.GetEnvOrDefault(testutil.EnvPostgresPass, "testpass")
	pgDB := testutil.GetEnvOrDefault(testutil.EnvPostgresDB, "integration_test")

	// Connect to PostgreSQL to create test table
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPass, pgAddr, pgDB)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create test table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS integration_test_table (
			column1 TEXT,
			column2 INTEGER
		)
	`
	_, err = conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Drop test table after test
	defer func() {
		dropTableSQL := `DROP TABLE IF EXISTS integration_test_table`
		_, dropErr := conn.Exec(context.Background(), dropTableSQL)
		if dropErr != nil {
			t.Logf("Warning: Failed to drop test table: %v", dropErr)
		}
	}()

	base := sink.BaseSink{
		ChanSize:      100,
		StreamName:    "",
		SinkAliasName: "PostgresSQL-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
	}
	testPostgresSQLConfig := config.PostgresSQLSinkConfig{
		Address:   pgAddr,
		Username:  pgUser,
		Password:  pgPass,
		Database:  pgDB,
		TableName: "integration_test_table",
		BulkSize:  5,
	}

	pgClient := NewPostgresSQLHandler(base, &testPostgresSQLConfig)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pgClient.WriteData()
	}()

	c := pgClient.GetFromTransformChan()
	for i := 1; i < 20; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				cast.ToString(i),
				i,
			},
			SinkName: "postgressql-1",
		}
	}

	// Wait for all data to be processed, then close sink and wait for goroutine exit
	time.Sleep(2 * time.Second)
	pgClient.CloseSink()
	wg.Wait()
}
