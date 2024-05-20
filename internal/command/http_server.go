// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package command

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gozix/di"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewHTTPServer is command constructor.
func NewHTTPServer(ctn di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "http-server",
		Short: "Run http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ctn.Call(func(e *echo.Echo, cfg *viper.Viper, logger *zap.Logger) (err error) {
				var addr = net.JoinHostPort(
					cfg.GetString("echo.host"),
					cfg.GetString("echo.port"),
				)

				logger.Info("Starting HTTP server", zap.String("addr", addr))

				go func() {
					runErr := e.Start(addr)
					if runErr == http.ErrServerClosed {
						logger.Info("Gracefully shutting down the HTTP server")
						return
					}
					if runErr != nil {
						logger.Error("HTTP server terminated incorrectly", zap.Error(runErr))
					}
				}()

				logger.Info("HTTP server started", zap.String("addr", addr))

				// wait, global context cancellation
				<-cmd.Context().Done()

				// graceful shutdown
				var timeout = 10 * time.Second
				logger.Info("Stopping HTTP server", zap.Duration("timeout", timeout))

				var ctx, cancel = context.WithTimeout(context.Background(), timeout)
				defer cancel()

				return e.Shutdown(ctx)
			})
		},
	}
}
