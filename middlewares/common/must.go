package common

import "github.com/gin-gonic/gin"

func Must(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			ResponseWriteError(c, err)
		}
	}
}
