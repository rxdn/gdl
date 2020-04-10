package request

import "errors"

type (
	ClientError error
	ServerError error
)

var (
	// 4xx
	ErrBadRequest      ClientError = errors.New("400 bad request")
	ErrUnauthorized    ClientError = errors.New("401 unauthorized")
	ErrForbidden       ClientError = errors.New("403 forbidden")
	ErrNotFound        ClientError = errors.New("404 not found")
	ErrTooManyRequests ClientError = errors.New("429 too many requests (ratelimit exceeded)")

	// 5xx
	ErrInternalServerError ServerError = errors.New("500 internal server error")
	ErrBadGateway          ServerError = errors.New("502 bad gateway")
	ErrServiceUnavailable  ServerError = errors.New("503 service unavailable")
	ErrGatewayTimeout      ServerError = errors.New("504 gateway timeout")

	ErrUnknown = errors.New("unknown error")

	errorCodes = map[int]error{
		400: ErrBadRequest,
		401: ErrUnauthorized,
		403: ErrForbidden,
		404: ErrNotFound,
		429: ErrTooManyRequests,
		500: ErrInternalServerError,
		502: ErrBadGateway,
		503: ErrServiceUnavailable,
		504: ErrGatewayTimeout,
	}
)

func IsServerError(err error) bool {
	_, ok := err.(ServerError)
	return ok
}

func IsClientError(err error) bool {
	_, ok := err.(ClientError)
	return ok
}
