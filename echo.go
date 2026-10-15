// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package echo

import (
	"github.com/gozix/di"
	"github.com/gozix/glue/v3"
	gzValidator "github.com/gozix/validator/v3"
	gzViper "github.com/gozix/viper/v3"
	gzZap "github.com/gozix/zap/v3"
	"github.com/labstack/echo/v4"

	"github.com/gozix/echo/v4/internal/command"
	"github.com/gozix/echo/v4/internal/configurator"
	"github.com/gozix/echo/v4/internal/factory"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct {
		svrNames []string
	}

	// Configurator is type alias of configurator.Configurator.
	Configurator = configurator.Configurator

	// Controller is type alias of controller.Controller.
	Controller = configurator.Controller

	// Factory is type alias of factory.Factory.
	Factory = factory.Factory
)

// BundleName is default definition name.
const BundleName = "echo"

var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle(svrNames ...string) *Bundle {
	return &Bundle{
		svrNames: svrNames,
	}
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder di.Builder) error {
	var opt = []di.BuilderOption{
		di.Provide(command.NewHTTPServer, glue.AsCliCommand()),
		di.Provide(factory.New, di.Constraint(4, withServerName(), di.Optional(true))),
	}

	for _, srvName := range b.svrNames {
		opt = append(opt, di.BuilderOptions(
			// server name
			di.Add(srvName, asServerName()),

			// echo
			di.Provide(
				func(f *Factory) (*echo.Echo, error) {
					return f.Echo(srvName)
				},
				asEcho(srvName),
			),
		))
	}

	return builder.Apply(opt...)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		gzValidator.BundleName,
		gzViper.BundleName,
		gzZap.BundleName,
	}
}
