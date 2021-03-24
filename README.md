## go-gin-extras
-----
#### Installation
-----

```bash
go get github.com/ms-xy/go-gin-extras
go get github.com/ms-xy/go-gin-extras@commithash
```

-----
#### Usage
-----

```go
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
    // do something with your session storage
    c.String(200, s.Token())
  })
  log.Fatal(engine.Run(environment.GetOrDefault("SERVICE_ADDRESS", "127.0.0.1:4000")))
}
```

-----
#### Available Middlewares
-----

- `middlewares/common.Logger()`

  Custom logger middleware that prints meaningful error messages and stack
  traces if attached via ctx.Set("error", err), see recovery middleware for an
  example.

- `middlewares/common.Recovery()`

  Custom recovery handler that gracefully recovers from panics, writes a 500
  response if possible and attaches detailed error information to the context
  for later retrieval and printing/analysis.

- `middlewares/session.DefaultSessionMiddleware()`

  Creates the default session wrapper using a mysql store with DSN taken from
  env (MYSQL_DATASOURCE). See SessionMiddleware(...) description below.

- `middlewares/session.SessionMiddleware(scs.Store)`

  Creates a thin wrapper around [scs](github.com/alexedwards/scs) using any
  supplied store as a session persistence backend. Refer to the scs
  documentation for indepth details and available store options.

  The following configuration options exist and can be freely changed before
  calling SessionMiddleware(...) in order to change the behavior of the created
  session manager:

  ```go
    // MySqlDataSource is the parameter used for creation of the store when calling DefaultSessionMiddleware()
    // See mysqlstore in scs for info on the required table schema
    MySqlDataSource string = env.GetOrDefault("MYSQL_DATASOURCE", "user:password@tcp(localhost)/databasename")

    // SessionCookie is the name of the session cookie used, defaults to 'session'
    SessionCookie string = env.GetOrDefault("SESSION_COOKIE", "session")
    // SessionDomain is the name of the domain associated with the session cookie
    SessionDomain string = env.GetOrDefault("SESSION_DOMAIN", "127.0.0.1")
    // SessionMaxAge is the maximum session lifetime in seconds
    SessionMaxAge int = mustParseInt(env.GetOrDefault("SESSION_MAX_AGE", "86400")) // seconds, 24 hour default
    // SessionIdleTimeout is the maximum idle time before a non-active session is discarded
    SessionIdleTimeout int = mustParseInt(env.GetOrDefault("SESSION_IDLE_TIMEOUT", "1800")) // seconds, 30 mins default
    // SessionSecure sets wether or not the cookie should be https only
    SessionSecure bool = strings.ToLower(os.Getenv("SESSION_SECURE")) == "true"
    // SessionHttpOnly sets wether the cookie be accessible via javascript
    SessionHttpOnly bool = (strings.ToLower(env.GetOrDefault("SESSION_HTTP_ONLY", "true")) == "true")
  ```

-----
#### License
-----

This project is licensed under the MIT license.
Please refer to the provided license file in the project root for details
