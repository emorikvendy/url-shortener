package resources

import (
	"database/sql"
	"fmt"

	"github.com/emorikvendy/url-shortener/internal/sources/datatypes"
	"github.com/emorikvendy/url-shortener/internal/sources/postgres"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

type R struct {
	Config   Config
	DB       *reform.DB
	conn     *sql.DB
	Adapters Adapters
}

type Config struct {
	DiagPort    int    `envconfig:"DIAG_PORT" default:"8081" required:"true"`
	RESTAPIPort int    `envconfig:"PORT" default:"8080" required:"true"`
	DBURL       string `envconfig:"DB_URL" default:"postgres://url_shortener:url_shortener@localhost:54320/url_shortener?sslmode=disable" required:"true"`
	Source      string `envconfig:"DATA_SOURCE" default:"postgres"`
	HashLen     int    `envconfig:"HASH_LEN" default:"10" required:"true"`
}

type Adapters struct {
	URL   datatypes.URLAdapter
	Stats datatypes.StatsAdapter
}

func New(logger *zap.SugaredLogger) (*R, error) {
	conf := Config{}
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, fmt.Errorf("can't process the config: %w", err)
	}
	if conf.HashLen < 8 || conf.HashLen > 32 {
		return nil, fmt.Errorf("can't process the config: hashLen must be between 8 and 32")
	}

	conn, err := sql.Open("pgx", conf.DBURL)
	if err != nil {
		return nil, err
	}

	db := reform.NewDB(conn, postgresql.Dialect, reform.NewPrintfLogger(logger.Infof))
	if conf.Source == "postgres" {
		urlAdapter := postgres.NewURLAdapter(db, conf.HashLen)
		stats := postgres.NewStatsAdapter(db)
		adapters := Adapters{
			URL:   urlAdapter,
			Stats: stats,
		}

		return &R{
			Config:   conf,
			DB:       db,
			Adapters: adapters,
		}, nil
	}

	return &R{
		Config: conf,
		DB:     db,
	}, nil
}

func (r *R) Release() error {
	return r.conn.Close()
}
