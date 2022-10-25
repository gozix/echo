// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import "github.com/labstack/echo/v4"

// Configurator is custom function configurator.
type Configurator func(*echo.Echo) error
