package common

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ms-xy/go-gin-extras/errors"
	color "gopkg.in/gookit/color.v1"
)

var (
	color2xx    = color.New(color.BgGreen, color.LightWhite)
	color4xx    = color.New(color.BgYellow, color.LightWhite)
	color5xx    = color.New(color.BgRed, color.LightWhite)
	colorElse   = color.New(color.BgWhite, color.Gray)
	colorMethod = color.New(color.BgBlue, color.LightWhite)
)

/*
Logger creates a gin.HandlerFunc that uses go's standard log package.
Use log.SetOutput to modify logging output destination.
Prefix is assembled by strings.Join(prefix, " ").
*/
func Logger(prefix ...string) gin.HandlerFunc {
	// define prefix - backwards compatibility by variadic argument
	var loggingPrefix string
	if prefix != nil || len(prefix) > 0 {
		loggingPrefix = strings.Join(prefix, " ")
	} else {
		loggingPrefix = "Server"
	}
	// define and return handler
	return func(c *gin.Context) {
		// get time for handling duration
		t := time.Now()

		// handle request / other middleware
		c.Next()

		// access the status we are sending
		// color code it appropriately
		status := errors.StatusCode(c.Writer.Status())
		var statusColor color.Style
		if status.Is2xxSuccess() {
			statusColor = color2xx
		} else if status.Is4xxClientError() {
			statusColor = color4xx
		} else if status.Is5xxServerError() {
			statusColor = color5xx
		} else {
			statusColor = colorElse
		}
		coloredStatus := statusColor.Sprintf("%3d", status.Int())

		// get request method and format appropriately
		coloredMethod := colorMethod.Sprintf("%6s", strings.ToUpper(c.Request.Method))

		// get duration
		latency := time.Since(t)

		// optional error, formatted (must start with a new line)
		errMsg := ""
		if err, exists := c.Get("error"); exists {
			if oErr, ok := err.(errors.Error); ok {
				if msg := oErr.Error(); msg != "" {
					errMsg = "\nPanic: " + msg
				}
				if data := oErr.Data(); data != nil {
					errMsg += fmt.Sprintf("\nAttached Data: %s", data)
				}
				if stackTrace := oErr.StackTrace(); stackTrace != "" {
					errMsg += "\nStack Trace:\n" + stackTrace
				}
				if errMsg != "" {
					errMsg = color.Red.Sprint(errMsg)
				}
			} else {
				errMsg = fmt.Sprintf("\nError: %s", err)
			}
		}

		log.Printf("[%s] %s [%s] %12s | %21s |%s %s -- %s\n",
			loggingPrefix,
			t.Format(time.RFC1123),
			coloredStatus,
			latency.String(),
			c.Request.RemoteAddr,
			coloredMethod,
			c.Request.RequestURI,
			errMsg,
		)
	}
}
