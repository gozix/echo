// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package command

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gozix/di"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gozix/echo/v4/internal/factory"
	"github.com/gozix/echo/v4/internal/filter"
)

// NewHTTPServer is command constructor.
func NewHTTPServer(ctn di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "http-server [name...]",
		Short: "Run http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ctn.Call(func(f *factory.Factory, cfg *viper.Viper, logger *zap.Logger) (err error) {
				var names []string
				if names, err = filter.Names(f.Names(), args); err != nil {
					return fmt.Errorf("unable filter echo servers : %w", err)
				}

				// build all servers before start
				var servers = make(map[string]*echo.Echo, len(names))
				for _, srvName := range names {
					if servers[srvName], err = f.Echo(srvName); err != nil {
						return err
					}
				}

				var wg, ctx = errgroup.WithContext(cmd.Context())

				for srvName, e := range servers {
					wg.Go(func() (err error) {
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

						var stopped = make(chan error, 1)
						go func() {
							stopped <- e.Start(addr)
						}()

						select {
						case <-ctx.Done():
							// graceful shutdown
							var timeout = 10 * time.Second
							log.Info("Stopping HTTP server", zap.Duration("timeout", timeout))

							var timeoutContext, cancel = context.WithTimeout(context.Background(), timeout)
							defer cancel()

							if err = e.Shutdown(timeoutContext); err != nil {
								return fmt.Errorf("unable to stop http server %s : %w", srvName, err)
							}

							err = <-stopped
						case err = <-stopped:
							// stopped without shutdown, e.g. the address is already in use
						}

						// http.ErrServerClosed is returned by the start after the shutdown
						if errors.Is(err, http.ErrServerClosed) {
							log.Info("HTTP server stopped")
							return nil
						}

						return fmt.Errorf("http server %s : %w", srvName, err)
					})
				}

				return wg.Wait()
			})
		},
	}
}
