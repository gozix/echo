package echo

import (
	"strconv"

	"github.com/gozix/di"

	"github.com/gozix/echo/v4/internal/tag"
)

// AsConfigurator is syntax sugar for the di container.
func AsConfigurator(srvName string) di.ProvideOption {
	return di.Tags{{
		Name: tag.Configurator + "." + srvName,
	}}
}

// AsController is syntax sugar for the di container.
func AsController(srvName string) di.ProvideOption {
	return di.ProvideOptions(di.Tags{{
		Name: tag.Controller + "." + srvName,
	}}, di.As(new(Controller)))
}

// AsMiddleware is syntax sugar for the di container.
func AsMiddleware(srvName string, priority int64) di.ProvideOption {
	return di.Tags{{
		Name: tag.Middleware + "." + srvName,
		Args: di.Args{{
			Key:   tag.ArgMiddlewarePriority,
			Value: strconv.FormatInt(priority, 10),
		}},
	}}
}

func asServerName() di.AddOption {
	return di.Tags{{
		Name: tag.ServerName,
	}}
}

func asEcho(srvName string) di.ProvideOption {
	return di.Tags{{
		Name: tag.Echo + "." + srvName,
	}}
}

func withServerName() di.Modifier {
	return di.WithTags(tag.ServerName)
}
