package util

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status codes
const (
	CodeSuccess        	= 200
	CodeValidationError	= 422
	CodeUnauthorized   	= 401
	CodeBadRequest     	= 400
	CodeForbidden      	= 403
	CodeNotAllowed     	= 405
	CodeUnknownError   	= 500
)

// APIResponse adalah struktur dasar untuk respons API
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewAPIResponse membangun respons API standar
func NewAPIResponse(status int, data interface{}, message string) APIResponse {
	if message == "" {
		message = buildDefaultMessage(status)
	}
	if data == nil {
		data = make(map[string]interface{})
	}
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// buildDefaultMessage memberikan pesan default berdasarkan status code
func buildDefaultMessage(status int) string {
	switch status {
	case CodeSuccess:
		return "Success"
	case CodeValidationError:
		return "Validation Error"
	case CodeUnauthorized:
		return "Unauthorized"
	case CodeBadRequest:
		return "Bad Request"
	case CodeForbidden:
		return "Forbidden"
	case CodeNotAllowed:
		return "Method Not Allowed"
	case CodeUnknownError:
		return "Internal Server Error"
	default:
		return "Unknown Error"
	}
}

// RespondWithError merespons kesalahan dengan logging tambahan
func RespondWithError(c *gin.Context, statusCode int, message string, data ...interface{}) {
	log.Printf("Error: %s, Status: %d", message, statusCode)
	if len(data) > 0 {
		c.JSON(statusCode, NewAPIResponse(statusCode, data[0], message))
	} else {
		c.JSON(statusCode, NewAPIResponse(statusCode, nil, message))
	}
}

// RespondWithSuccess merespons sukses
func RespondWithSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, NewAPIResponse(CodeSuccess, data, message))
}
