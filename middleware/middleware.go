package middleware

import (
	"echoapp/container"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasttemplate"

	"gopkg.in/boj/redistore.v1"
)

// InitLoggerMiddleware initialize a middleware for logger.
func InitLoggerMiddleware(e *echo.Echo, container container.Container) {
	e.Use(RequestLoggerMiddleware(container))
	e.Use(ActionLoggerMiddleware(container))
}

// InitSessionMiddleware initialize a middleware for session management.
func InitSessionMiddleware(e *echo.Echo, container container.Container) {
	conf := container.GetConfig()
	logger := container.GetLogger()

	if conf.Extension.SecurityEnabled {
		if conf.Redis.Enabled {
			logger.GetZapLogger().Infof("Try redis connection")
			address := fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port)
			store, err := redistore.NewRediStore(conf.Redis.ConnectionPoolSize, "tcp", address, "", []byte("secret"))
			if err != nil {
				logger.GetZapLogger().Errorf("Failure redis connection")
			}
			e.Use(session.Middleware(store))
			logger.GetZapLogger().Infof(fmt.Sprintf("Success redis connection, %s", address))
		} else {
			e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
		}
	}
}

// RequestLoggerMiddleware is middleware for logging the contents of requests.
func RequestLoggerMiddleware(container container.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			if err := next(c); err != nil {
				c.Error(err)
			}

			template := fasttemplate.New(container.GetConfig().Log.RequestLogFormat, "${", "}")
			logstr := template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "remote_ip":
					return w.Write([]byte(c.RealIP()))
				case "uri":
					return w.Write([]byte(req.RequestURI))
				case "method":
					return w.Write([]byte(req.Method))
				case "status":
					return w.Write([]byte(strconv.Itoa(res.Status)))
				default:
					return w.Write([]byte(""))
				}
			})
			container.GetLogger().GetZapLogger().Infof(logstr)
			return nil
		}
	}
}

// ActionLoggerMiddleware is middleware for logging the start and end of controller processes.
// ref: https://echo.labstack.com/cookbook/middleware
func ActionLoggerMiddleware(container container.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := container.GetLogger()
			logger.GetZapLogger().Debugf(c.Path() + " Action Start")
			if err := next(c); err != nil {
				c.Error(err)
			}
			logger.GetZapLogger().Debugf(c.Path() + " Action End")
			return nil
		}
	}
}

// equalPath judges whether a given path contains in the path list.
func equalPath(cpath string, paths []string) bool {
	for i := range paths {
		if regexp.MustCompile(paths[i]).Match([]byte(cpath)) {
			return true
		}
	}
	return false
}
