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

	"github.com/gozix/echo/v4/internal/command"
	"github.com/gozix/echo/v4/internal/configurator"
	"github.com/gozix/echo/v4/internal/echo"
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
	}

	for _, srvName := range b.svrNames {
		opt = append(opt, di.BuilderOptions(
			// server name
			di.Add(srvName, asServerName(srvName)),

			// echo
			di.Provide(
				echo.New, asEcho(srvName),
				di.Constraint(0, withConfigurator(srvName)),
			),

			// configurators
			di.Provide(
				configurator.NewController, AsConfigurator(srvName),
				di.Constraint(0, withController(srvName), di.Optional(true)),
			),
			di.Provide(
				configurator.NewMiddleware, AsConfigurator(srvName),
				di.Constraint(0, withMiddleware(srvName), di.Optional(true), sortByPriority()),
			),
			di.Provide(
				configurator.NewEcho, AsConfigurator(srvName),
				di.Constraint(0, withServerName(srvName)),
			),
			di.Provide(configurator.NewErrHandler, AsConfigurator(srvName)),
			di.Provide(configurator.NewValidator, AsConfigurator(srvName)),
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
