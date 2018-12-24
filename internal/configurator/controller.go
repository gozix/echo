// Package configurator provides dependency injection definitions.
package configurator

import (
	"github.com/labstack/echo"
	"github.com/sarulabs/di"
)

// Controller is a controller interface.
type Controller interface {
	Serve(e *echo.Echo)
}

const (
	// DefControllerConfiguratorName is a definition name.
	DefControllerConfiguratorName = "echo.configurator.controller"

	// TagController is a controller tag name.
	TagController = "echo.controller"
)

// DefControllerConfigurator is echo controller definition getter.
func DefControllerConfigurator() di.Def {
	return di.Def{
		Name: DefControllerConfiguratorName,
		Tags: []di.Tag{{
			Name: TagConfigurator,
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return func(e *echo.Echo) (err error) {
				for name, def := range ctn.Definitions() {
					for _, tag := range def.Tags {
						if tag.Name != TagController {
							continue
						}

						var c Controller
						if err = ctn.Fill(name, &c); err != nil {
							return err
						}

						c.Serve(e)

						break
					}
				}

				return nil
			}, nil
		},
	}
}
