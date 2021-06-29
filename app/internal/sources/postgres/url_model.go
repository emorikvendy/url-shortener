package postgres

import (
	"github.com/emorikvendy/url-shortener/internal/datatypes"
	"time"
)

//go:generate reform

//reform:url
type urlModel struct {
	ID        int64     `reform:"id,pk"`
	Name      string    `reform:"name"`
	Link      string    `reform:"link"`
	Hash      string    `reform:"hash"`
	CreatedAt time.Time `reform:"created_at"`
	UpdatedAt time.Time `reform:"updated_at"`
}

// BeforeInsert set CreatedAt and UpdatedAt.
func (model *urlModel) BeforeInsert() error {
	model.CreatedAt = time.Now().UTC().Truncate(time.Second)
	model.UpdatedAt = model.CreatedAt
	return nil
}

// BeforeUpdate set UpdatedAt.
func (model *urlModel) BeforeUpdate() error {
	model.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	return nil
}

func (model *urlModel) toBaseUrl() *datatypes.URL {
	external := datatypes.URL{
		ID:   &model.ID,
		Name: &model.Name,
		Hash: &model.Hash,
		Link: &model.Link,
	}

	return &external
}
