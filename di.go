package echo

import (
	"strconv"

	"github.com/gozix/di"
)

const (
	// tagConfigurator is tag to mark middleware's.
	tagConfigurator = "echo.configurator"

	// tagMiddleware is tag to mark middleware's.
	tagMiddleware = "echo.middleware"

	// argMiddlewarePriority is name of priority argument.
	argMiddlewarePriority = "priority"
)

// AsConfigurator is syntax sugar for the di container.
func AsConfigurator() di.ProvideOption {
	return di.Tags{{
		Name: tagConfigurator,
	}}
}

// AsController is syntax sugar for the di container.
func AsController() di.ProvideOption {
	return di.As(new(Controller))
}

// AsMiddleware is syntax sugar for the di container.
func AsMiddleware(priority int64) di.ProvideOption {
	return di.Tags{{
		Name: tagMiddleware,
		Args: di.Args{{
			Key:   argMiddlewarePriority,
			Value: strconv.FormatInt(priority, 10),
		}},
	}}
}

func sortByPriority() di.Modifier {
	return di.Sort(func(x, y di.Definition) bool {
		var xp int64
		for _, tag := range x.Tags() {
			for _, arg := range tag.Args {
				if arg.Key == argMiddlewarePriority {
					xp, _ = strconv.ParseInt(arg.Value, 10, 64)
				}
			}
		}

		var yp int64
		for _, tag := range y.Tags() {
			for _, arg := range tag.Args {
				if arg.Key == argMiddlewarePriority {
					yp, _ = strconv.ParseInt(arg.Value, 10, 64)
				}
			}
		}

		return xp > yp
	})
}

func withConfigurator() di.Modifier {
	return di.WithTags(tagConfigurator)
}

func withMiddleware() di.Modifier {
	return di.WithTags(tagMiddleware)
}
