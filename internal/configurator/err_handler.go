// Package configurator provides dependency injection definitions.
package configurator

import (
	"net/http"

	"github.com/gozix/glue/v2"
	zapBundle "github.com/gozix/zap/v2"
	"github.com/labstack/echo"
	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"

	"github.com/gozix/echo/v2/errors"
)

// DefErrHandlerConfiguratorName definition name.
const DefErrHandlerConfiguratorName = "echo.configurator.err_handler"

// DefErrHandlerConfigurator is echo error handler definition getter.
func DefErrHandlerConfigurator() di.Def {
	return di.Def{
		Name: DefErrHandlerConfiguratorName,
		Tags: []di.Tag{{
			Name: TagConfigurator,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var registry glue.Registry
			if err = ctn.Fill(glue.DefRegistry, &registry); err != nil {
				return nil, err
			}

			var defName string
			if err = registry.Fill(DefErrHandlerConfiguratorName, &defName); err != nil {
				return nil, err
			}

			if len(defName) > 0 {
				var handler Configurator
				if err = ctn.Fill(defName, &handler); err != nil {
					return nil, err
				}

				return handler, nil
			}

			var logger *zap.Logger
			if err = ctn.Fill(zapBundle.BundleName, &logger); err != nil {
				return nil, err
			}

			return func(e *echo.Echo) (err error) {
				e.HTTPErrorHandler = func(err error, c echo.Context) {
					var (
						e      = c.Echo()
						code   = http.StatusInternalServerError
						msg    = http.StatusText(code)
						errMsg interface{}
					)

					switch he := err.(type) {
					case *echo.HTTPError:
						code = he.Code
						msg = he.Message.(string)
					case *errors.HTTPError:
						code = he.Code
						msg = he.Message
						errMsg = he.Description
					}

					if e.Debug {
						msg = err.Error()
					}

					if !c.Response().Committed {
						if c.Request().Method == echo.HEAD {
							err = c.NoContent(code)
						} else {
							var m = echo.Map{"message": msg}
							if errMsg != nil {
								m["error"] = errMsg
							}
							err = c.JSON(code, m)
						}
						if err != nil {
							logger.Error("An error occurred", zap.Error(err))
						}
					}
				}

				return nil
			}, nil
		},
	}
}
