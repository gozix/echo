// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/gozix/echo/v3/errors"
)

// NewErrHandler is configurator constructor.
func NewErrHandler(logger *zap.Logger) Configurator {
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
	}
}
