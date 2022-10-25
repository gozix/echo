// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Decorator implementation.
type Decorator struct {
	validator *validator.Validate
}

// NewValidator is echo configurator constructor.
func NewValidator(v *validator.Validate) Configurator {
	return func(e *echo.Echo) error {
		e.Validator = &Decorator{
			validator: v,
		}

		return nil
	}
}

// Validate data
func (v *Decorator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
