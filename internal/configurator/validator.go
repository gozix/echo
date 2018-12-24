// Package configurator provides dependency injection definitions.
package configurator

import (
	validatorBundle "github.com/gozix/validator"
	"github.com/labstack/echo"
	"github.com/sarulabs/di"
	"gopkg.in/go-playground/validator.v9"
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
				var v *validatorBundle.Validate
				if err = ctn.Fill(validatorBundle.BundleName, &v); err != nil {
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
