// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package middleware

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gozix/echo/v3/errors"
)

// ErrTransConfig defines the config for ErrTransWithConfig middleware.
type ErrTransConfig struct {
	Skipper    middleware.Skipper
	Translator *ut.UniversalTranslator
}

// ErrTransWithConfig returns echo.MiddlewareFunc.
func ErrTransWithConfig(conf ErrTransConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if conf.Skipper != nil && conf.Skipper(c) {
				return next(c)
			}

			if err = next(c); err == nil {
				return nil
			}

			if httpErr, ok := err.(*errors.HTTPError); ok {
				if e, ok := httpErr.Cause.(validator.ValidationErrors); ok {
					var trans = make(validator.ValidationErrorsTranslations)
					for i := range e {
						var fe = e[i]
						trans[fe.Field()] = fe.Translate(conf.Translator.GetFallback())
					}
					httpErr.Description = trans
					return httpErr
				}
			}

			return err
		}
	}
}
