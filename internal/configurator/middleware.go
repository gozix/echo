// Package configurator provides dependency injection definitions.
package configurator

import (
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
)

// Middleware is middleware definition.
type Middleware struct {
	Func     echo.MiddlewareFunc
	Priority int
}

const (
	// DefMiddlewareConfiguratorName is a definition name.
	DefMiddlewareConfiguratorName = "echo.configurator.middleware"

	// TagMiddleware is tag to mark middleware's.
	TagMiddleware = "echo.middleware"
)

// DefMiddlewareConfigurator is echo middleware definition getter.
func DefMiddlewareConfigurator() di.Def {
	return di.Def{
		Name: DefMiddlewareConfiguratorName,
		Tags: []di.Tag{{
			Name: TagConfigurator,
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return func(e *echo.Echo) (err error) {
				// extract middleware's
				var mws = make([]*Middleware, 0, 4)
				for key, def := range ctn.Definitions() {
					for _, tag := range def.Tags {
						if tag.Name != TagMiddleware {
							continue
						}

						var mw = new(Middleware)
						if v, ok := tag.Args["priority"]; ok {
							if mw.Priority, err = strconv.Atoi(v); err != nil {
								return err
							}
						}

						if err = ctn.Fill(key, &mw.Func); err != nil {
							return err
						}

						mws = append(mws, mw)
						break
					}
				}

				// sort middleware's
				sort.Slice(mws, func(i, j int) bool {
					return mws[i].Priority > mws[j].Priority
				})

				// register middleware's
				for _, mw := range mws {
					e.Use(mw.Func)
				}

				return nil
			}, nil
		},
	}
}
