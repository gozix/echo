// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package echo

import (
	"github.com/labstack/echo/v4"

	"github.com/gozix/echo/v4/internal/configurator"
)

func New(configurators []configurator.Configurator) (*echo.Echo, error) {
	var e = echo.New()
	for _, c := range configurators {
		if err := c(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}
