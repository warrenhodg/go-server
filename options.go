package server

import (
	"time"

	"github.com/warrenhodg/health"
	"go.uber.org/zap"
)

type Options struct {
	addr             string
	logger           *zap.Logger
	health           health.IHealth
	shutdownDuration time.Duration
	warningDuration  time.Duration
}

func DefaultOptions() *Options {
	return &Options{}
}

// WithHealth sets the health checker used to allow health checks
// to fail before actually refusing to accept new connections
func (o *Options) WithHealth(value health.IHealth) *Options {
	o.health = value
	return o
}

// WithListenAddress sets the address on which the server will listen
func (o *Options) WithListenAddress(value string) *Options {
	o.addr = value
	return o
}

// WithLogger sets the logger for the server
func (o *Options) WithLogger(value *zap.Logger) *Options {
	o.logger = value.With(zap.String("module", "server"))
	return o
}

// WithWarningDuration sets the duration for which the health
// check should fail before closing the tcp listener. This allows
// time for load balancers to remove this instance before it stops
// actually refusing connections
func (o *Options) WithWarningDuration(value time.Duration) *Options {
	o.warningDuration = value
	return o
}

// WithShutdownDuration sets the duration for which existing connections
// will be allowed to safely terminate before they are focefully closed
func (o *Options) WithShutdownDuration(value time.Duration) *Options {
	o.shutdownDuration = value
	return o
}
