# server

A golang gin-gonic compatible server, supporting safe shutdowns with health checks.

## Usage

```
import (
    "context"
	"github.com/warrenhodg/health"
	"github.com/warrenhodg/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	h := health.New(logger)

	svrOptions := server.
		DefaultOptions().
		WithListenAddress(":8080").
		WithHealth(h).
		WithLogger(logger).
		WithWarningDuration(time.Second * 5)
		WithShutdownDuration(time.Second * 5)

	svr, err := server.New(svrOptions)
    if err != nil {
        panic(err)
    }

    defer svr.Close()

    h.RegisterEndpoint(svr)

    svr.GET("Hello", func (ctx *gin.Context) {
        ctx.String("200", "Hello")
    })

    // Wait until app must end
}
```

## Notes

On shutdown, this server first starts failing health checks, then stops accepting new connections but waits for existing connections to safely
terminate, and finally closes any existing connections forcefully.