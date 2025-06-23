// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// NewEcho is echo configurator constructor.
func NewEcho(srvName string, cfg *viper.Viper) Configurator {
	return func(e *echo.Echo) error {
		if !cfg.IsSet("echo." + srvName) {
			return fmt.Errorf("configuration echo.%s is not found", srvName)
		}

		var c = cfg.Sub("echo." + srvName)

		switch c.GetString("level") {
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

		if c.IsSet("static") {
			e.Static(
				c.GetString("static.prefix"),
				c.GetString("static.root"),
			)
		}

		e.Debug = c.GetBool("debug")
		e.HidePort = c.GetBool("hide_port")
		e.HideBanner = c.GetBool("hide_banner")

		return nil
	}
}
