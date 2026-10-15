// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package tag

const (
	// ServerName is tag to mark a server name value.
	ServerName = "echo.server_name"

	// Echo is tag to mark a echo instance.
	Echo = "echo.echo"

	// Configurator is tag to mark configurators.
	Configurator = "echo.configurator"

	// Controller is tag to mark controllers.
	Controller = "echo.controller"

	// Middleware is tag to mark middlewares.
	Middleware = "echo.middleware"

	// ArgMiddlewarePriority is name of priority argument.
	ArgMiddlewarePriority = "priority"
)
