package main

import (
	"github.com/ms-xy/go-common/environment"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ms-xy/go-gin-extras/middlewares/common"
	"github.com/ms-xy/go-gin-extras/middlewares/session"
)

func main() {
	engine := gin.New()
	log.SetOutput(io.MultiWriter(os.Stdout))
	engine.Use(common.Logger())
	engine.Use(common.Recovery())
	engine.Use(session.DefaultSessionMiddleware())
	engine.GET("/", func(c *gin.Context) {
		s := session.GetSession(c)
		c.String(200, s.Token())
	})
	log.Fatal(engine.Run(environment.GetOrDefault("SERVICE_ADDRESS", "127.0.0.1:4000")))
}
