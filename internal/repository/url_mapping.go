package repository

import (
	"github.com/url-shortner/internal/model"
	"gorm.io/gorm"
)

type urlMappingRepo struct {
	db *gorm.DB
}

func NewUrlMappingRepo(db *gorm.DB) URLMappingRepo {
	return &urlMappingRepo{
		db: db,
	}
}

func (r *urlMappingRepo) Create(urlMapping *model.URLMapping) error {
	return r.db.Create(urlMapping).Error
}

func (r *urlMappingRepo) GetByShortCode(shortCode string) (*model.URLMapping, error) {
	var urlMapping model.URLMapping
	err := r.db.Where("short_code = ?", shortCode).First(&urlMapping).Error
	return &urlMapping, err
}

func (r *urlMappingRepo) GetByLongURL(longURL string) (*model.URLMapping, error) {
	var urlMapping model.URLMapping
	err := r.db.Where("long_url = ?", longURL).First(&urlMapping).Error
	return &urlMapping, err
}

func (r *urlMappingRepo) Update(urlMapping *model.URLMapping) error {
	return r.db.Save(urlMapping).Error
}

func (r *urlMappingRepo) UpdateClickAndAccess(id string) error {
	return r.db.
		Model(&model.URLMapping{}).
		Where("id = ?", id).
		Update("clicked_count", gorm.Expr("clicked_count + ?", 1)).
		Update("last_accessed_at", gorm.Expr("now()")).
		Error
}

func (r *urlMappingRepo) List(query model.Pagination) ([]model.URLMapping, int, error) {
	var urlMappings []model.URLMapping

	var count int64
	err := r.db.Model(&model.URLMapping{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Model(&model.URLMapping{}).Limit(query.Limit).Offset((query.Offset - 1) * query.Limit).Find(&urlMappings).Error
	return urlMappings, int(count), err
}
