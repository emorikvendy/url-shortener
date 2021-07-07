package postgres

import (
	"github.com/emorikvendy/url-shortener/internal/datatypes"
)

//go:generate reform

//reform:stats
type statsModel struct {
	URLID int64 `reform:"url_id,pk"`
	Hits  int64 `reform:"hits"`
}

func (model *statsModel) toBaseStats() *datatypes.Stats {
	external := datatypes.Stats{
		URLID: &model.URLID,
		Hits:  &model.Hits,
	}

	return &external
}
