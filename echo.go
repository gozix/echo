// Package echo provides dependency injection definitions.
package echo

import (
	"github.com/gozix/glue/v2"
	validatorBundle "github.com/gozix/validator/v2"
	viperBundle "github.com/gozix/viper/v2"
	zapBundle "github.com/gozix/zap/v2"
	"github.com/labstack/echo"
	"github.com/sarulabs/di/v2"

	"github.com/gozix/echo/v2/internal/command"
	"github.com/gozix/echo/v2/internal/configurator"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct {
		errHandlerDefName string
	}

	// Configurator is type alias of configurator.Configurator.
	Configurator = configurator.Configurator

	// Controller is type alias of controller.Controller.
	Controller = configurator.Controller

	// Option interface.
	Option interface {
		apply(b *Bundle)
	}

	// optionFunc wraps a func so it satisfies the Option interface.
	optionFunc func(b *Bundle)
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

// ErrHandler option.
func ErrHandler(defName string) Option {
	return optionFunc(func(b *Bundle) {
		b.errHandlerDefName = defName
	})
}

// NewBundle create bundle instance.
func NewBundle(options ...Option) *Bundle {
	var b = new(Bundle)

	for _, option := range options {
		option.apply(b)
	}

	return b
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
				var registry glue.Registry
				if err = ctn.Fill(glue.DefRegistry, &registry); err != nil {
					return nil, err
				}

				registry.Set(configurator.DefErrHandlerConfiguratorName, b.errHandlerDefName)

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

// apply implements Option.
func (f optionFunc) apply(bundle *Bundle) {
	f(bundle)
}
