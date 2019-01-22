// Package middleware provide implementations of custom middleware for the echo framework.
package middleware

import (
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"

	"github.com/gozix/echo/errors"
)

// ZapConfig defines the config for AccessWithConfig middleware.
type ZapConfig struct {
	// Logged fields.
	//
	// - id (Request ID)
	// - remote_ip
	// - host
	// - method
	// - uri
	// - path
	// - protocol
	// - referer
	// - user_agent
	// - status
	// - error
	// - latency (In nanoseconds)
	// - latency_human (Human readable)
	// - bytes_in (Bytes received)
	// - bytes_out (Bytes sent)
	// - header:<NAME>
	// - query:<NAME>
	// - form:<NAME>
	//
	// Optional. Default value DefaultZapConfig.Fields.
	Fields []string

	// Skipper defines a function to skip middleware.
	Logger *zap.Logger

	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// DefaultZapConfig is the default Zap middleware config.
var DefaultZapConfig = ZapConfig{
	Fields: []string{
		"id", "remote_ip", "host", "method", "uri", "user_agent", "status",
		"error", "latency", "latency_human", "bytes_in", "bytes_out",
	},
}

// ZapWithConfig returns echo.MiddlewareFunc.
func ZapWithConfig(cfg ZapConfig) echo.MiddlewareFunc {
	if len(cfg.Fields) == 0 {
		cfg.Fields = DefaultZapConfig.Fields
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if cfg.Skipper != nil && cfg.Skipper(c) {
				return next(c)
			}

			var (
				start   = time.Now()
				logFunc = cfg.Logger.Info
			)

			if err = next(c); err != nil {
				logFunc = cfg.Logger.Error
			}

			var (
				req    = c.Request()
				res    = c.Response()
				fields = make([]zap.Field, 0, len(cfg.Fields))
				finish = time.Now()
			)

			for _, field := range cfg.Fields {
				switch field {
				case "id":
					var id = req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}

					fields = append(fields, zap.String(field, id))
				case "remote_ip":
					fields = append(fields, zap.String(field, c.RealIP()))
				case "host":
					fields = append(fields, zap.String(field, req.Host))
				case "method":
					fields = append(fields, zap.String(field, req.Method))
				case "uri":
					fields = append(fields, zap.String(field, req.RequestURI))
				case "path":
					var p = req.URL.Path
					if p == "" {
						p = "/"
					}

					fields = append(fields, zap.String(field, p))
				case "protocol":
					fields = append(fields, zap.String(field, req.Proto))
				case "referer":
					fields = append(fields, zap.String(field, req.Referer()))
				case "user_agent":
					fields = append(fields, zap.String(field, req.UserAgent()))
				case "status":
					var s = res.Status
					if err != nil {
						switch v := err.(type) {
						case *errors.HTTPError:
							s = v.Code
						case *echo.HTTPError:
							s = v.Code
						default:
							s = echo.ErrInternalServerError.Code
						}
					}

					fields = append(fields, zap.Int(field, s))
				case "error":
					if err != nil {
						fields = append(fields, zap.Error(err))
					}
				case "latency":
					fields = append(fields, zap.Duration(field, finish.Sub(start)))
				case "latency_human":
					fields = append(fields, zap.String(field, finish.Sub(start).String()))
				case "bytes_in":
					var cl = req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}

					fields = append(fields, zap.String(field, cl))
				case "bytes_out":
					fields = append(fields, zap.Int64(field, res.Size))
				default:
					switch {
					case strings.HasPrefix(field, "header:"):
						fields = append(fields, zap.String(field, c.Request().Header.Get(field[7:])))
					case strings.HasPrefix(field, "query:"):
						fields = append(fields, zap.String(field, c.QueryParam(field[6:])))
					case strings.HasPrefix(field, "form:"):
						fields = append(fields, zap.String(field, c.FormValue(field[5:])))
					case strings.HasPrefix(field, "cookie:"):
						if cookie, err := c.Cookie(field[7:]); err == nil {
							fields = append(fields, zap.String(field, cookie.Value))
						}
					}
				}
			}

			logFunc(req.Method+" "+req.RequestURI, fields...)

			return err
		}
	}
}
