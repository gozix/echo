// Package configurator provide container definitions.
package configurator

import "github.com/labstack/echo"

// Configurator is custom function configurator.
type Configurator = func(*echo.Echo) error

// TagConfigurator is a configurator tag name.
const TagConfigurator = "echo.configurator"
