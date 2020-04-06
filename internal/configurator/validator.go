// Package configurator provides dependency injection definitions.
package configurator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"

	gzValidator "github.com/gozix/validator/v2"
)

// Wrapper implementation.
type Wrapper struct {
	validator *validator.Validate
}

// DefValidatorConfiguratorName is a definition name.
const DefValidatorConfiguratorName = "echo.configurator.validator"

// DefValidatorConfigurator is echo validator definition getter.
func DefValidatorConfigurator() di.Def {
	return di.Def{
		Name: DefValidatorConfiguratorName,
		Tags: []di.Tag{{
			Name: TagConfigurator,
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return func(e *echo.Echo) (err error) {
				var v *validator.Validate
				if err = ctn.Fill(gzValidator.BundleName, &v); err != nil {
					return err
				}

				e.Validator = &Wrapper{
					validator: v,
				}

				return nil
			}, nil
		},
	}
}

// Validate data
func (v *Wrapper) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
