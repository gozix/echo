// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// NewEcho is echo configurator constructor.
func NewEcho(cfg *viper.Viper) Configurator {
	return func(e *echo.Echo) error {
		if cfg.IsSet("echo.debug") {
			e.Debug = cfg.GetBool("echo.debug")
		}

		switch cfg.GetString("echo.level") {
		case "debug":
			e.Logger.SetLevel(log.DEBUG)
		case "info":
			e.Logger.SetLevel(log.INFO)
		case "warn":
			e.Logger.SetLevel(log.WARN)
		case "error":
			e.Logger.SetLevel(log.ERROR)
		case "off":
			e.Logger.SetLevel(log.OFF)
		}

		if cfg.IsSet("echo.static") {
			e.Static(
				cfg.GetString("echo.static.prefix"),
				cfg.GetString("echo.static.root"),
			)
		}

		if cfg.IsSet("echo.hide_banner") {
			e.HideBanner = cfg.GetBool("echo.hide_banner")
		}

		if cfg.IsSet("echo.hide_port") {
			e.HidePort = cfg.GetBool("echo.hide_port")
		}

		return nil
	}
}
