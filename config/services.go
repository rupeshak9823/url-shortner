package config

import "github.com/url-shortner/internal/service/url_service"

type ServiceList struct {
	URLService *url_service.URLService
}
