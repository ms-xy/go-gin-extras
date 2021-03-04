package common

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/ms-xy/go-gin-extras/errors"
	"github.com/ms-xy/go-gin-extras/stack"
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

	// distinguish between error types
	switch statusCode := oError.StatusCode(); {
	case statusCode.Is501NotImplemented():
		c.JSON(501, gin.H{"error": `501 Not Implemented`})
	case statusCode.Is4xxClientError():
		c.JSON(oError.StatusCode().Int(), gin.H{"error": oError.Error()})
	default:
		c.JSON(500, gin.H{"error": `500 Internal Server Error`})
	}

	// set error so logging handler can attach it
	c.Set("error", oError)

	// prevent further handlers
	c.Abort()
}
