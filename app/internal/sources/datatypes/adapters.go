package datatypes

import dt "github.com/emorikvendy/url-shortener/internal/datatypes"

type URLAdapter interface {
	GetByHash(hash string) (*dt.URL, error)
	GetByID(id int64) (*dt.URL, error)
	URLExists(url string) (bool, error)
	HashExists(hash string) (bool, error)
	Store(*dt.URL) error
	GetList(params map[string]interface{}) ([]dt.URL, error)
	Delete(id int64) error
}

type StatsAdapter interface {
	GetByURLID(id int64) (*dt.Stats, error)
	AddByURLID(id int64) error
}
