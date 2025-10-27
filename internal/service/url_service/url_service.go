package url_service

import (
	"context"
	"fmt"
	"time"

	redisCache "github.com/go-redis/cache/v9"
	"github.com/url-shortner/internal/constant"
	"github.com/url-shortner/internal/model"
	"github.com/url-shortner/internal/repository"
	"github.com/url-shortner/internal/util"
	"gorm.io/gorm"
)

type URLService interface {
	Shorten(urlMappingCreate URLMappingCreate) (*model.URLMapping, error)
	GetByShortCode(shortCode string) (*model.URLMapping, error)
	List(query URLMappingListInput) ([]model.URLMapping, int, error)
}

type urlService struct {
	urlMappingRepo repository.URLMappingRepo
	redisCache     *redisCache.Cache
}

func NewUrlService(urlMappingRepo repository.URLMappingRepo, redisCache *redisCache.Cache) URLService {
	return &urlService{
		urlMappingRepo: urlMappingRepo,
		redisCache:     redisCache,
	}
}

func (s urlService) GetByShortCode(shortCode string) (*model.URLMapping, error) {
	ctx := context.Background()

	urlMapping := &model.URLMapping{}
	err := s.redisCache.Get(ctx, shortCode, urlMapping)
	if err == redisCache.ErrCacheMiss {
		// Not in cache, fetch from DB
		urlMapping, err = s.urlMappingRepo.GetByShortCode(shortCode)
		if err != nil {
			return nil, err
		}

		// Store in cache
		err = s.redisCache.Set(&redisCache.Item{
			Key:   shortCode,
			Value: urlMapping,
			TTL:   2000 * time.Minute,
		})
		if err != nil {
			fmt.Println("Redis set error:", err)
		}
	} else if err != nil {
		return nil, err
	}

	go func(urlMapping model.URLMapping) {
		err := s.urlMappingRepo.UpdateClickAndAccess(urlMapping.ID)
		if err != nil {
			fmt.Println("UpdateClickAndAccess error:", err)
		}
	}(*urlMapping)

	return urlMapping, nil
}

func (s urlService) Shorten(urlMappingCreate URLMappingCreate) (*model.URLMapping, error) {
	longURL := urlMappingCreate.LongURL
	if longURL == "" || !util.IsValidWebURL(longURL) {
		return nil, util.BadRequestError(util.BadRequestError{
			Message: "Long url is not valid",
		})
	}
	mapping, err := s.urlMappingRepo.GetByLongURL(longURL)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		return mapping, nil
	}
	suffix := 0
	length := constant.URLShortLength
	var shortUrl, inputUrl = "", longURL
	for {
		shortUrl = util.StringToBase62(inputUrl, length)
		_, err := s.urlMappingRepo.GetByShortCode(shortUrl)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err == gorm.ErrRecordNotFound {
			break
		}
		suffix++
		if suffix > 1000 {
			shortUrl = util.StringToBase62(inputUrl, length+1)
			break
		}
		inputUrl = fmt.Sprintf("%s#%d", longURL, suffix)
	}

	urlMapping := model.URLMapping{
		ShortCode: shortUrl,
		LongURL:   longURL,
	}
	err = s.urlMappingRepo.Create(&urlMapping)
	if err != nil {
		return nil, err
	}
	return &urlMapping, nil
}

func (s urlService) List(query URLMappingListInput) ([]model.URLMapping, int, error) {
	if query.PerPage == 0 {
		query.PerPage = 10
	}
	return s.urlMappingRepo.List(model.Pagination{
		Limit:  query.PerPage,
		Offset: query.Page,
	})
}
