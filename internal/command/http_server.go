// Package command contains cli command definitions.
package command

import (
	"context"
	"net"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	gzGlue "github.com/gozix/glue/v2"
	gzViper "github.com/gozix/viper/v2"
	gzZap "github.com/gozix/zap/v2"
)

// DefEchoHTTPServerName is command definition name.
const DefEchoHTTPServerName = "cli.cmd.echo.http-server"

// DefEchoHTTPServer is command definition getter.
func DefEchoHTTPServer(bundleName string) di.Def {
	return di.Def{
		Name: DefEchoHTTPServerName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return &cobra.Command{
				Use:   "http-server",
				Short: "Run http server",
				RunE: func(cmd *cobra.Command, args []string) (err error) {
					var e *echo.Echo
					if err = ctn.Fill(bundleName, &e); err != nil {
						return err
					}

					var cfg *viper.Viper
					if err = ctn.Fill(gzViper.BundleName, &cfg); err != nil {
						return err
					}

					var logger *zap.Logger
					if err = ctn.Fill(gzZap.BundleName, &logger); err != nil {
						return err
					}

					var addr = net.JoinHostPort(
						cfg.GetString("echo.host"),
						cfg.GetString("echo.port"),
					)

					logger.Info("Starting HTTP server", zap.String("addr", addr))

					go func() {
						if err = e.Start(addr); err != nil {
							logger.Info("Gracefully shutting down the HTTP server")
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
				},
			}, nil
		},
	}
}
