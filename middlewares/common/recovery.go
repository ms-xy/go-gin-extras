package common

import (
	"go-gin-extras/errors"
	"go-gin-extras/stack"

	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return DeferredGracefulPanicRecovery
}

func DeferredGracefulPanicRecovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ResponseWriteError(c, r)
		}
	}()
	c.Next()
}

func ResponseWriteError(c *gin.Context, _err interface{}) {
	var oError errors.Error
	if err, ok := _err.(errors.Error); ok {
		oError = err
	} else if str, ok := _err.(string); ok {
		oError = errors.NewError(500, str, stack.Stack(3), nil)
	} else if err, ok := _err.(error); ok {
		oError = errors.NewError(500, err.Error(), stack.Stack(3), nil)
	} else {
		oError = errors.NewError(
			500, `Unexpected Error Type`,
			string(debug.Stack()), gin.H{"_err": _err})
	}

	// do not disclose internal errors to the client
	// detailed error info is only provided for malformed client requests
	statusCode := oError.StatusCode()
	if statusCode.Is501NotImplemented() {
		c.JSON(501, gin.H{"error": `501 Not Implemented`})
	} else if statusCode.Is4xxClientError() {
		c.JSON(statusCode.Int(), gin.H{"error": oError.Error()})
	} else {
		c.JSON(500, gin.H{"error": `500 Internal Server Error`})
	}

	// set error so logging handler can attach it
	c.Set("error", oError)

	// prevent subsequent handlers (add logging before recovery!)
	c.Abort()
}
