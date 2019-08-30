// Package middleware provide implementations of custom middleware for the echo framework.
package middleware

import (
	"github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gozix/echo/v2/errors"
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
