package util

import (
	"context"
	"github.com/labstack/echo/v4"
)

type ContextLogger struct {
	Logger echo.Logger
}

func Logger(ctx context.Context) echo.Logger {
	logger := ctx.Value("logger").(*ContextLogger)
	return (*logger).Logger
}

func GetUserFromContext(ctx context.Context) *int64 {
	id := ctx.Value("USER_ID").(*int64)

	if id != nil {
		return id
	}

	return nil
}
