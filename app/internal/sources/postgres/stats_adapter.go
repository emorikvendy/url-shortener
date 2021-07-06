package postgres

import (
	"fmt"

	dt "github.com/emorikvendy/url-shortener/internal/datatypes"

	"gopkg.in/reform.v1"
)

type StatsAdapter struct {
	db *reform.DB
}

func NewStatsAdapter(db *reform.DB) *StatsAdapter {
	return &StatsAdapter{
		db: db,
	}
}

func (adapter StatsAdapter) GetByURLID(id int64) (*dt.Stats, error) {
	rec, err := adapter.db.FindByPrimaryKeyFrom(statsModelTable, id)
	if err != nil {
		return nil, fmt.Errorf("couldn't find stats by url_id: %w", err)
	}

	internal := rec.(*statsModel)

	return internal.toBaseStats(), nil

}

func (adapter StatsAdapter) AddByURLID(id int64) error {
	_, err := adapter.db.Query("UPDATE stats SET hits=hits+1 WHERE url_id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't update stats by url_id: %w", err)
	}

	return nil
}
