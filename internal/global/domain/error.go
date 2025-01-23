package dm

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/erx/erx"
	"gorm.io/gorm"
)

var (
	ErrNotImplement     = errors.New("not implemented")
	ErrEmptyRequest     = errors.New("empty request body")
	ErrUnknown          = errors.New("unknown")
	ErrLoginFailed      = errors.New("login failed")
	ErrUserNameRepetite = errors.New("user name repetite")
)

// Convert original error to http status code
func ToHttpCode(err error) int {
	ogErr := erx.RUnwrap(err)
	switch ogErr {
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
	case ErrNotImplement:
		return http.StatusNotImplemented
	case ErrEmptyRequest:
		return http.StatusBadRequest
	case ErrUnknown:
		return http.StatusInternalServerError
	case ErrLoginFailed:
		return http.StatusUnauthorized
	case ErrUserNameRepetite:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func RespErrWithCode(c *gin.Context, err error) {
	statusCode := ToHttpCode(err)
	RespErr(c, statusCode, err)
}
