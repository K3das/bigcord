package archiver

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/K3das/bigcord/scraping/store"
	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Archiver struct {
	log *zap.SugaredLogger

	discord *discordgo.Session

	store     *store.Store
	warehouse driver.Conn
}

func NewArchiver(ctx context.Context, rawLog *zap.Logger, store *store.Store, warehouse driver.Conn, token string) (*Archiver, error) {
	s := &Archiver{}
	s.log = rawLog.Sugar().With("source", "archiver")

	account, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("error creating discordgo session: %s", err)
	}
	account.Client = &http.Client{
		Timeout: 20 * time.Second,
		Transport: promhttp.InstrumentRoundTripperCounter(discordRequestsCounter,
			promhttp.InstrumentRoundTripperDuration(discordRequestDurationHistogram,
				http.DefaultTransport,
			),
		),
	}

	s.discord = account

	user, err := s.discord.User("@me", discordgo.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error creating discordgo session: %s", err)
	}

	s.log = s.log.With("discord_user_id", user.ID)
	s.log.Infof("account works")

	s.store = store
	s.warehouse = warehouse

	return s, nil
}
