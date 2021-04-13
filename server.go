package server

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IServer interface {
	gin.IRouter
	io.Closer
}

type Server struct {
	*gin.Engine

	options  *Options
	listener net.Listener
}

func New(options *Options) (svr *Server, err error) {
	engine := gin.New()

	listener, err := net.Listen("tcp", options.addr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Engine:   engine,
		options:  options,
		listener: listener,
	}

	if options.logger != nil {
		options.logger.Info("running", zap.String("module", "server"), zap.String("addr", options.addr))
	}

	if options.health != nil {
		server.options.health.SetSystemState("server", true)
	}

	go engine.RunListener(listener)

	return server, nil
}

func (svr *Server) Close() error {
	if svr.options.logger != nil {
		svr.options.logger.Info("shutdown notice")
	}

	if svr.options.health != nil {
		svr.options.health.SetSystemState("server", false)
	}

	if svr.options.warningDuration > time.Duration(0) {
		time.Sleep(svr.options.warningDuration)
	}

	if svr.options.logger != nil {
		svr.options.logger.Info("shutdown start")
	}

	ctx := context.Background()
	if svr.options.shutdownDuration > time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, svr.options.shutdownDuration)
		defer cancel()
	}

	return svr.Engine.Shutdown(ctx)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
