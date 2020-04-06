// Package command contains cli command definitions.
package command

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
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
					var cnf *viper.Viper
					if err = ctn.Fill(viper.BundleName, &cnf); err != nil {
						return err
					}

					var e *echo.Echo
					if err = ctn.Fill("echo", &e); err != nil {
						return err
					}

					go e.Logger.Error(
						e.Start(
							net.JoinHostPort(
								cnf.GetString("echo.host"),
								cnf.GetString("echo.port"),
							),
						),
					)

					signalChan := make(chan os.Signal, 1)
					signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

					<-signalChan

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					return e.Shutdown(ctx)
				},
			}, nil
		},
	}
}
