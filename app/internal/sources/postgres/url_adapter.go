package postgres

import (
	"crypto/md5"
	"fmt"
	dt "github.com/emorikvendy/url-shortener/internal/datatypes"
	"gopkg.in/reform.v1"
	"strings"
)

type URLAdapter struct {
	db      *reform.DB
	hashLen int
}

func New(db *reform.DB, hashLen int) *URLAdapter {
	return &URLAdapter{
		db:      db,
		hashLen: hashLen,
	}
}
func (adapter URLAdapter) GetByHash(hash string) (*dt.URL, error) {
	args := make([]interface{}, 1)
	tail := "WHERE hash = $1"
	args[0] = hash

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return nil, err
	} else if len(sts) > 1 {
		return nil, fmt.Errorf("found more than one link with hash %s", hash)
	} else if len(sts) == 0 {
		return nil, nil
	}
	internal := sts[0].(*urlModel)

	return internal.toBaseUrl(), nil
}

func (adapter URLAdapter) GetByID(id int64) (*dt.URL, error) {
	args := make([]interface{}, 1)
	tail := "WHERE id = $1"
	args[0] = id

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return nil, err
	} else if len(sts) == 0 {
		return nil, nil
	}
	internal := sts[0].(*urlModel)

	return internal.toBaseUrl(), nil

}

func (adapter URLAdapter) getByID(id int64) (*urlModel, error) {
	args := make([]interface{}, 1)
	tail := "WHERE id = $1"
	args[0] = id

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return nil, err
	} else if len(sts) == 0 {
		return nil, nil
	}
	internal := sts[0].(*urlModel)

	return internal, nil

}

func (adapter URLAdapter) Delete(id int64) error {
	args := make([]interface{}, 1)
	tail := "WHERE id = $1"
	args[0] = id

	_, err := adapter.db.DeleteFrom(urlModelTable, tail, args...)
	if err != nil {
		return err
	}

	return nil

}

func (adapter URLAdapter) URLExists(url string) (bool, error) {

	args := make([]interface{}, 1)
	tail := "WHERE link = $1 LIMIT 1"
	args[0] = url

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return false, err
	} else if len(sts) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
func (adapter URLAdapter) HashExists(hash string) (bool, error) {
	args := make([]interface{}, 1)
	tail := "WHERE hash = $1 LIMIT 1"
	args[0] = hash

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return false, err
	} else if len(sts) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (adapter URLAdapter) Store(URL *dt.URL) error {
	var (
		internal *urlModel
		err      error
	)
	if URL.ID != nil && *URL.ID > 0 {
		internal, err = adapter.getByID(*URL.ID)
		if err != nil {
			return err
		}
	}
	if internal == nil {
		internal = &urlModel{
			Name: *URL.Name,
			Link: *URL.Link,
		}
		link := internal.Link
		for i := 0; i < 100; i++ {
			hs := fmt.Sprintf("%x", md5.Sum([]byte(link)))[0:adapter.hashLen]
			if ok, err2 := adapter.HashExists(hs); ok {
				link = link + hs
				if i == 100 {
					return fmt.Errorf("can't get hash for link %s", internal.Link)
				}
				continue
			} else if err2 != nil {
				return err2
			}
			internal.Hash = hs
			URL.Hash = &hs
			break
		}
	} else {
		if internal.Link != *URL.Link {
			return fmt.Errorf("it is forbidden to change the link of an existing url")
		}
		internal.Name = *URL.Name
	}

	err = adapter.db.Save(internal)
	if err != nil {
		return err
	}

	URL.ID = &internal.ID
	return nil
}

func (adapter URLAdapter) GetList(params map[string]interface{}) ([]dt.URL, error) {

	intersect := make([]string, len(params))
	args := make([]interface{}, len(params))

	i := 0
	for key, arg := range params {
		intersect[i] = fmt.Sprintf("%s = $%d", key, i+1)
		args[i] = arg
		i++
	}

	tail := "WHERE " + strings.Join(intersect, " AND ")

	sts, err := adapter.db.SelectAllFrom(urlModelTable, tail, args...)
	if err != nil {
		return nil, err
	} else if len(sts) == 0 {
		return nil, nil
	}

	result := make([]dt.URL, len(sts), len(sts))

	for i := 0; i < len(sts); i++ {
		internal := sts[i].(*urlModel)
		result[i] = *internal.toBaseUrl()
	}

	return result, nil
}
