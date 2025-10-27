package util

import (
	"net/http"

	"github.com/url-shortner/internal/model"
)

type ServiceError struct {
	Code    string
	Message string
}

func (f ServiceError) Error() string {
	if f.Message != "" {
		return f.Message
	}
	return f.Code
}

type ForbiddenError ServiceError

func (f ForbiddenError) Error() string {
	return ServiceError(f).Error()
}

type Unauthorized ServiceError

func (f Unauthorized) Error() string {
	return ServiceError(f).Error()
}

type NotFoundError ServiceError

func (f NotFoundError) Error() string {
	return ServiceError(f).Error()
}

type BadRequestError ServiceError

func (f BadRequestError) Error() string {
	return ServiceError(f).Error()
}

type InternalServerError ServiceError

func (f InternalServerError) Error() string {
	return ServiceError(f).Error()
}

func ConvertErrorToHttpError(err error) model.HTTPError {
	switch err := err.(type) {
	case BadRequestError:
		return model.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Message,
		}
	default:
		return model.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
}

func ReturnHttpError(err error, w http.ResponseWriter) {
	httpError := ConvertErrorToHttpError(err)
	w.WriteHeader(httpError.Code)
	w.Write([]byte(httpError.Message))
}
