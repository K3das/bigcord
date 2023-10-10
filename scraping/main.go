package main

import (
	"context"
	"fmt"
	"github.com/K3das/bigcord/scraping/api"
	"github.com/K3das/bigcord/scraping/archiver"
	"github.com/K3das/bigcord/scraping/jobs"
	"github.com/K3das/bigcord/scraping/store"
	"github.com/K3das/bigcord/scraping/warehouse"
	"github.com/caarlos0/env/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var CommitHash = "INVALID_COMMIT"

type config struct {
	MetricsHost string `env:"METRICS_HOST,required" envDefault:"0.0.0.0:9090"`
	WebHost     string `env:"WEB_HOST,required" envDefault:"0.0.0.0:3926"`

	SQLiteDSN string `env:"SQLITE_DSN,required" envDefault:":memory:"`

	ClickhouseHosts    []string `env:"CLICKHOUSE_HOSTS,required"`
	ClickhouseDatabase string   `env:"CLICKHOUSE_DATABASE,required"`
	ClickhouseUsername string   `env:"CLICKHOUSE_USERNAME,required"`
	ClickhousePassword string   `env:"CLICKHOUSE_PASSWORD,required"`

	DiscordToken string `env:"DISCORD_TOKEN,required"`

	ChannelIDs []string `env:"CHANNEL_IDS,required"`
}

func main() {
	// TODO: pls disregard my shameless panicking - I'm normal, I swear
	cfg := config{}
	if err := env.ParseWithOptions(&cfg, env.Options{
		Prefix: "SCRAPING_",
	}); err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	rawLog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}

	rawLog = rawLog.With(zap.Field{
		Key:    "commit",
		Type:   zapcore.StringType,
		String: CommitHash,
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.MetricsHost, nil); err != nil {
			panic(fmt.Errorf("metrics server: %w", err))
		}
	}()

	log := rawLog.Sugar().With("source", "main")
	log.Info("starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqliteStore, err := store.NewStore(ctx, rawLog, cfg.SQLiteDSN)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create store: %w", err))
	}

	clickhouse, err := warehouse.NewWarehouseConnection(ctx, rawLog, cfg.ClickhouseHosts, cfg.ClickhouseDatabase, cfg.ClickhouseUsername, cfg.ClickhousePassword)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create warehouse: %w", err))
	}

	a, err := archiver.NewArchiver(ctx, rawLog, sqliteStore, clickhouse, cfg.DiscordToken)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create archiver: %w", err))
	}

	jobStore := jobs.NewJobStore(ctx, rawLog)

	web := api.NewAPI(rawLog, a, sqliteStore, jobStore)

	err = web.Listen(ctx, cfg.WebHost)
	if err != nil {
		panic(fmt.Errorf("error serving on %s: %w", cfg.WebHost, err))
	}
}
