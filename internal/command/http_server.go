// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package command

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gozix/di"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gozix/echo/v4/internal/modifier"
)

const (
	tagEcho       = "echo.echo"
	tagServerName = "echo.server_name"
)

// NewHTTPServer is command constructor.
func NewHTTPServer(ctn di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "http-server [name...]",
		Short: "Run http server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var modServers *modifier.Modifier
			if modServers, err = modifier.NewModifier(tagServerName, args); err != nil {
				return fmt.Errorf("unable create echo glob modifier : %w", err)
			}

			return ctn.Call(func(serverNames []string, cfg *viper.Viper, logger *zap.Logger) error {
				var wg, ctx = errgroup.WithContext(cmd.Context())

				for _, srvName := range serverNames {
					var srvName = srvName
					wg.Go(func() (err error) {
						var e *echo.Echo
						if err = ctn.Resolve(&e, di.WithTags(tagEcho+"."+srvName)); err != nil {
							return err
						}

						var (
							subCfg = cfg.Sub("echo." + srvName)
							addr   = net.JoinHostPort(
								subCfg.GetString("host"),
								subCfg.GetString("port"),
							)
							log = logger.With(
								zap.String("name", srvName),
								zap.String("addr", addr),
							)
						)

						log.Info("Starting HTTP server")

						go func() {
							if err = e.Start(addr); err != nil {
								log.Info("Gracefully shutting down the HTTP server")
							}
						}()

						log.Info("HTTP server started")

						// wait, context cancellation
						<-ctx.Done()

						// graceful shutdown
						var timeout = 10 * time.Second
						log.Info("Stopping HTTP server", zap.Duration("timeout", timeout))

						var timeoutContext, cancel = context.WithTimeout(context.Background(), timeout)
						defer cancel()

						return e.Shutdown(timeoutContext)
					})
				}

				return wg.Wait()
			}, di.Constraint(0, modServers.Modifier()))
		},
	}
}
