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

	"github.com/gozix/echo/v3/internal/command"
	"github.com/gozix/echo/v3/internal/configurator"
	"github.com/gozix/echo/v3/internal/echo"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Configurator is type alias of configurator.Configurator.
	Configurator = configurator.Configurator

	// Controller is type alias of controller.Controller.
	Controller = configurator.Controller
)

// BundleName is default definition name.
const BundleName = "echo"

var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder di.Builder) error {
	return builder.Apply(
		// echo
		di.Provide(echo.New, di.Constraint(0, di.Optional(true), withConfigurator())),

		// command's
		di.Provide(command.NewHTTPServer, glue.AsCliCommand()),

		// configurator's
		di.Provide(configurator.NewController, AsConfigurator()),
		di.Provide(configurator.NewEcho, AsConfigurator()),
		di.Provide(configurator.NewErrHandler, AsConfigurator()),
		di.Provide(
			configurator.NewMiddleware, AsConfigurator(),
			di.Constraint(0, withMiddleware(), sortByPriority()),
		),
		di.Provide(configurator.NewValidator, AsConfigurator()),
	)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		gzValidator.BundleName,
		gzViper.BundleName,
		gzZap.BundleName,
	}
}
