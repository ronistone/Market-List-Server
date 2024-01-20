package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/ronistone/market-list/src/util"
)

func InjectLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := c.Request()
		c.SetRequest(request.WithContext(context.WithValue(ctx, "logger", &util.ContextLogger{Logger: c.Logger()})))
		return next(c)
	}
}

func InjectUserId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var userId int64 = 1
		ctx := c.Request().Context()
		request := c.Request()
		c.SetRequest(request.WithContext(context.WithValue(ctx, "USER_ID", &userId)))
		return next(c)
	}
}
