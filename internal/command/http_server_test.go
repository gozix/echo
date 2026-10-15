// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package command_test

import (
	"context"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gozix/di"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gozix/echo/v4/internal/command"
	"github.com/gozix/echo/v4/internal/factory"
	"github.com/gozix/echo/v4/internal/tag"
)

func TestHTTPServerGracefulShutdown(t *testing.T) {
	var ctn = newContainer(t, "0")

	var f *factory.Factory
	require.NoError(t, ctn.Resolve(&f))

	var e, err = f.Echo("admin")
	require.NoError(t, err)

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var done = make(chan error, 1)
	go func() {
		done <- execute(ctx, ctn)
	}()

	require.Eventually(t, func() bool {
		return e.ListenerAddr() != nil
	}, 5*time.Second, 10*time.Millisecond, "http server is not started")

	resp, err := http.Get("http://" + e.ListenerAddr().String())
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	cancel()

	select {
	case err = <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("http server is not stopped")
	}
}

func TestHTTPServerStartError(t *testing.T) {
	var busy, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	defer func() {
		_ = busy.Close()
	}()

	_, port, err := net.SplitHostPort(busy.Addr().String())
	require.NoError(t, err)

	var done = make(chan error, 1)
	go func() {
		done <- execute(context.Background(), newContainer(t, port))
	}()

	select {
	case err = <-done:
		require.ErrorIs(t, err, syscall.EADDRINUSE)
	case <-time.After(5 * time.Second):
		t.Fatal("http server start error is not returned")
	}
}

func newContainer(t *testing.T, port string) di.Container {
	t.Helper()

	var cfg = viper.New()
	cfg.Set("echo.admin.host", "127.0.0.1")
	cfg.Set("echo.admin.port", port)
	cfg.Set("echo.admin.hide_banner", true)
	cfg.Set("echo.admin.hide_port", true)

	var builder, err = di.NewBuilder(
		di.Add(cfg),
		di.Add(zap.NewNop()),
		di.Add(validator.New()),
		di.Provide(
			factory.New,
			di.Constraint(4, di.WithTags(tag.ServerName), di.Optional(true)),
		),
	)
	require.NoError(t, err)

	ctn, err := builder.Build()
	require.NoError(t, err)

	return ctn
}

func execute(ctx context.Context, ctn di.Container) error {
	var cmd = command.NewHTTPServer(ctn)
	cmd.SetArgs([]string{})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	return cmd.ExecuteContext(ctx)
}
