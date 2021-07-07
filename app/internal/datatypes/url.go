package datatypes

type URL struct {
	ID   *int64  `json:"id"`
	Name *string `json:"name"`
	Link *string `json:"link"`
	Hash *string `json:"hash"`
}
