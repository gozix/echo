package echo

import (
	"strconv"

	"github.com/gozix/di"
)

const (
	// tagServerName is tag to mark a server name value.
	tagServerName = "echo.server_name"

	// tagEcho is tag to mark a echo instance.
	tagEcho = "echo.echo"

	// tagConfigurator is tag to mark configurators.
	tagConfigurator = "echo.configurator"

	// tagController is tag to mark controllers.
	tagController = "echo.controller"

	// tagMiddleware is tag to mark middlewares.
	tagMiddleware = "echo.middleware"

	// argMiddlewarePriority is name of priority argument.
	argMiddlewarePriority = "priority"
)

// AsConfigurator is syntax sugar for the di container.
func AsConfigurator(srvName string) di.ProvideOption {
	return di.Tags{{
		Name: tagConfigurator + "." + srvName,
	}}
}

// AsController is syntax sugar for the di container.
func AsController(srvName string) di.ProvideOption {
	return di.ProvideOptions(di.Tags{{
		Name: tagController + "." + srvName,
	}}, di.As(new(Controller)))
}

// AsMiddleware is syntax sugar for the di container.
func AsMiddleware(srvName string, priority int64) di.ProvideOption {
	return di.Tags{{
		Name: tagMiddleware + "." + srvName,
		Args: di.Args{{
			Key:   argMiddlewarePriority,
			Value: strconv.FormatInt(priority, 10),
		}},
	}}
}

func asServerName() di.AddOption {
	return di.Tags{{
		Name: tagServerName,
	}}
}

func asEcho(srvName string) di.ProvideOption {
	return di.Tags{{
		Name: tagEcho + "." + srvName,
	}}
}

func withServerName() di.Modifier {
	return di.WithTags(tagServerName)
}
