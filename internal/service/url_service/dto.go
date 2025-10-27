package url_service

type URLMappingCreate struct {
	LongURL string `json:"longURL"`
}

type URLMappingListInput struct {
	Page    int
	PerPage int
}
