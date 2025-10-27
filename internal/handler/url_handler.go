package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/url-shortner/internal/service/url_service"
	"github.com/url-shortner/internal/util"
)

type URLHandler struct {
	urlService url_service.URLService
}

func NewURLHandler(urlService url_service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h URLHandler) Shorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req url_service.URLMappingCreate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			httpError := util.ConvertErrorToHttpError(err)
			w.WriteHeader(httpError.Code)
			w.Write([]byte(httpError.Message))
			return
		}
		urlMapping, err := h.urlService.Shorten(req)
		if err != nil {
			httpError := util.ConvertErrorToHttpError(err)
			w.WriteHeader(httpError.Code)
			w.Write([]byte(httpError.Message))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(urlMapping.ShortCode))
	}
}

func (h URLHandler) GetByShortCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "shortCode")
		urlMapping, err := h.urlService.GetByShortCode(shortCode)
		if err != nil {
			httpError := util.ConvertErrorToHttpError(err)
			w.WriteHeader(httpError.Code)
			w.Write([]byte(httpError.Message))
			return
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(urlMapping.LongURL))
	}
}

func (h URLHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req url_service.URLMappingListInput
		req.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		req.PerPage, _ = strconv.Atoi(r.URL.Query().Get("perPage"))
		urlMappings, total, err := h.urlService.List(req)
		if err != nil {
			util.ReturnHttpError(err, w)
			return
		}
		output, err := json.Marshal(map[string]interface{}{
			"total": total,
			"data":  urlMappings,
		})
		if err != nil {
			util.ReturnHttpError(err, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(output))
	}
}
