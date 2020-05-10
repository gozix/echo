// Package configurator provides dependency injection definitions.
package configurator

import (
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
)

// Controller is a controller interface.
type Controller interface {
	Serve(e *echo.Echo)
}

// GracefulController is a graceful interface.
type GracefulController interface {
	OnShutdown()
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
						switch v := c.(type) {
						case GracefulController:
							e.Server.RegisterOnShutdown(v.OnShutdown)
						}

						break
					}
				}

				return nil
			}, nil
		},
	}
}
