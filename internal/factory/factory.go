// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package factory

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gozix/di"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gozix/echo/v4/internal/configurator"
	internalEcho "github.com/gozix/echo/v4/internal/echo"
	"github.com/gozix/echo/v4/internal/tag"
)

// configKey is config section with servers configuration.
const configKey = "echo"

// Factory builds echo servers on demand, after config is read.
type Factory struct {
	ctn      di.Container
	cfg      *viper.Viper
	logger   *zap.Logger
	validate *validator.Validate
	names    []string

	mux     sync.Mutex
	servers map[string]*echo.Echo
}

// New is factory constructor. The names are the server names passed to the bundle, may be empty.
func New(ctn di.Container, cfg *viper.Viper, logger *zap.Logger, validate *validator.Validate, names []string) *Factory {
	return &Factory{
		ctn:      ctn,
		cfg:      cfg,
		logger:   logger,
		validate: validate,
		names:    names,
		servers:  make(map[string]*echo.Echo, len(names)),
	}
}

// Names returns the server names passed to the bundle or, if none were passed, the keys of the echo config section.
func (f *Factory) Names() []string {
	if len(f.names) > 0 {
		return slices.Clone(f.names)
	}

	return slices.Sorted(maps.Keys(f.cfg.GetStringMap(configKey)))
}

// Echo returns the server echo instance. The instance is built on the first call.
func (f *Factory) Echo(name string) (_ *echo.Echo, err error) {
	f.mux.Lock()
	defer f.mux.Unlock()

	if e, ok := f.servers[name]; ok {
		return e, nil
	}

	if !f.cfg.IsSet(configKey + "." + name) {
		return nil, fmt.Errorf("configuration %s.%s is not found", configKey, name)
	}

	var configurators []configurator.Configurator
	if configurators, err = f.getConfigurators(name); err != nil {
		return nil, fmt.Errorf("unable to resolve echo %s configurators : %w", name, err)
	}

	var e *echo.Echo
	if e, err = internalEcho.New(configurators); err != nil {
		return nil, err
	}

	f.servers[name] = e

	return e, nil
}

func (f *Factory) getConfigurators(name string) (_ []configurator.Configurator, err error) {
	var configurators []configurator.Configurator
	err = f.ctn.Call(
		func(controllers []configurator.Controller, mws []echo.MiddlewareFunc, custom []configurator.Configurator) {
			configurators = append([]configurator.Configurator{
				configurator.NewController(controllers),
				configurator.NewMiddleware(mws),
				configurator.NewEcho(name, f.cfg),
				configurator.NewErrHandler(f.logger),
				configurator.NewValidator(f.validate),
			}, custom...)
		},
		di.Constraint(0, di.Optional(true), di.WithTags(tag.Controller+"."+name)),
		di.Constraint(1, di.Optional(true), di.WithTags(tag.Middleware+"."+name), sortByPriority()),
		di.Constraint(2, di.Optional(true), di.WithTags(tag.Configurator+"."+name)),
	)

	return configurators, err
}

func sortByPriority() di.Modifier {
	return di.Sort(func(x, y di.Definition) bool {
		var xp int64
		for _, t := range x.Tags() {
			for _, arg := range t.Args {
				if arg.Key == tag.ArgMiddlewarePriority {
					xp, _ = strconv.ParseInt(arg.Value, 10, 64)
				}
			}
		}

		var yp int64
		for _, t := range y.Tags() {
			for _, arg := range t.Args {
				if arg.Key == tag.ArgMiddlewarePriority {
					yp, _ = strconv.ParseInt(arg.Value, 10, 64)
				}
			}
		}

		return xp > yp
	})
}
