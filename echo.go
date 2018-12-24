// Package echo provides dependency injection definitions.
package echo

import (
	validatorBundle "github.com/gozix/validator"
	viperBundle "github.com/gozix/viper"
	zapBundle "github.com/gozix/zap"
	"github.com/labstack/echo"
	"github.com/sarulabs/di"

	"github.com/gozix/echo/internal/command"
	"github.com/gozix/echo/internal/configurator"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Configurator is type alias of configurator.Configurator.
	Configurator = configurator.Configurator

	// Controller is type alias of controller.Controller.
	Controller = configurator.Controller
)

const (
	// TagController is alias of controller.Controller.
	TagController = configurator.TagController

	// TagConfigurator is alias of configurator.TagConfigurator.
	TagConfigurator = configurator.TagConfigurator

	// TagMiddleware is alias of configurator.TagMiddleware.
	TagMiddleware = configurator.TagMiddleware
)

// BundleName is default definition name.
const BundleName = "echo"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(
		// echo
		di.Def{
			Name: BundleName,
			Build: func(ctn di.Container) (_ interface{}, err error) {
				var e = echo.New()
				for name, def := range ctn.Definitions() {
					for _, tag := range def.Tags {
						if tag.Name != configurator.TagConfigurator {
							continue
						}

						var conf configurator.Configurator
						if err = ctn.Fill(name, &conf); err != nil {
							return nil, err
						}

						if err = conf(e); err != nil {
							return nil, err
						}

						break
					}
				}

				return e, nil
			},
		},

		// command's
		command.DefEchoHTTPServer(),

		// configurator's
		configurator.DefControllerConfigurator(),
		configurator.DefEchoConfigurator(),
		configurator.DefErrHandlerConfigurator(),
		configurator.DefMiddlewareConfigurator(),
		configurator.DefValidatorConfigurator(),
	)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		viperBundle.BundleName,
		validatorBundle.BundleName,
		zapBundle.BundleName,
	}
}
