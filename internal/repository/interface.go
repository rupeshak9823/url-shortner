package repository

import "github.com/url-shortner/internal/model"

type URLMappingRepo interface {
	Create(urlMapping *model.URLMapping) error
	GetByShortCode(shortCode string) (*model.URLMapping, error)
	GetByLongURL(longURL string) (*model.URLMapping, error)
	Update(urlMapping *model.URLMapping) error
	UpdateClickAndAccess(id string) error
	List(query model.Pagination) ([]model.URLMapping, int, error)
}
