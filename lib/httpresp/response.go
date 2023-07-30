package httpresp

import (
	"net/http"

	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/gin-gonic/gin"
)

const (
	SuccessCode = 0
	ErrorCode   = 1
)

// Send a normal success response with JSON data
func SendData(c *gin.Context, data interface{}, statusCode int) {
	c.JSON(statusCode, StandardDataResponse{
		Code:    SuccessCode,
		Message: "",
		Data:    data,
	})
}

// Send an error JSON response based on custom error (see error package)
func SendError(c *gin.Context, err error) {
	customErr, ok := err.(resperror.CustomErrWithCode)
	if !ok {
		c.JSON(http.StatusInternalServerError, StandardResponse{
			Code:    resperror.UnknownError,
			Message: "",
		})
		return
	}
	c.JSON(customErr.StatusCode, StandardResponse{
		Code:    customErr.Code,
		Message: customErr.Error(),
	})
}

// If err is not a recognized error, fallback to fallbackErr
func SendErrorWithFallback(c *gin.Context, err error, fallbackErr resperror.CustomErrWithCode) {
	customErr, ok := err.(resperror.CustomErrWithCode)
	if !ok {
		c.JSON(fallbackErr.StatusCode, StandardResponse{
			Code:    fallbackErr.Code,
			Message: fallbackErr.Message,
		})
		return
	}
	c.JSON(customErr.StatusCode, StandardResponse{
		Code:    customErr.Code,
		Message: customErr.Error(),
	})
}

// Avoid use (use SendError instead)
func SendErrorWithStatusCode(c *gin.Context, err error, statusCode int) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.JSON(statusCode, StandardResponse{
		Code:    ErrorCode,
		Message: errMsg,
	})
}

// Send error message (Avoid use, use SendError instead)
func SendErrorMsg(c *gin.Context, errMsg string, statusCode int) {
	c.JSON(statusCode, StandardResponse{
		Code:    ErrorCode,
		Message: errMsg,
	})
}

// Send normal success response
func SendSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, StandardResponse{
		Code:    SuccessCode,
		Message: "",
	})
}
