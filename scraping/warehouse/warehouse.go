package warehouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
)

func NewWarehouseConnection(ctx context.Context, rawLog *zap.Logger, hosts []string, database, username, password string) (driver.Conn, error) {
	log := rawLog.Sugar().With("source", "warehouse")

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: hosts,
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug:  true,
		Debugf: log.Debugf,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		//Compression: &clickhouse.Compression{
		//	Method: clickhouse.CompressionLZ4,
		//},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "scraping", Version: "0.1"},
			},
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         50,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating clickhouse connection: %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging clickhouse: %w", err)
	}

	log.Info("connected to clickhouse")

	return conn, nil
}
