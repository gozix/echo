// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"github.com/labstack/echo/v4"
)

// Controller is a controller interface.
type Controller interface {
	Serve(e *echo.Echo)
}

// NewController is echo configurator constructor.
func NewController(controllers []Controller) Configurator {
	return func(e *echo.Echo) error {
		for _, controller := range controllers {
			controller.Serve(e)
		}

		return nil
	}
}
