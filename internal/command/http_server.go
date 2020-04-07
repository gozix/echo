// Package command contains cli command definitions.
package command

import (
	"context"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/gozix/glue/v2"
	"github.com/gozix/viper/v2"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
)

// DefEchoHTTPServerName is command definition name.
const DefEchoHTTPServerName = "cli.cmd.echo.http-server"

// DefEchoHTTPServer is command definition getter.
func DefEchoHTTPServer() di.Def {
	return di.Def{
		Name: DefEchoHTTPServerName,
		Tags: []di.Tag{{
			Name: glue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return &cobra.Command{
				Use:   "http-server",
				Short: "Run http server",
				RunE: func(cmd *cobra.Command, args []string) (err error) {
					var cfg *viper.Viper
					if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
						return err
					}

					var e *echo.Echo
					if err = ctn.Fill("echo", &e); err != nil {
						return err
					}

					// run
					go e.Logger.Error(
						e.Start(
							net.JoinHostPort(
								cfg.GetString("echo.host"),
								cfg.GetString("echo.port"),
							),
						),
					)

					// wait
					<-cmd.Context().Done()

					// shutdown
					var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					return e.Shutdown(ctx)
				},
			}, nil
		},
	}
}
