// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"github.com/labstack/echo/v4"
)

// NewMiddleware is echo configurator constructor.
func NewMiddleware(mws []echo.MiddlewareFunc) Configurator {
	return func(e *echo.Echo) error {
		for _, mw := range mws {
			e.Use(mw)
		}

		return nil
	}
}
