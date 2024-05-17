package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

func (app *Application) DeleteOldFileIfExists(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.FileID != "" {

			if err := app.Client.DeleteFile(context.Background(), app.FileID); err != nil {
				c.Logger().Error(err)
				return internalServerError(c)
			}

			if app.ThreadID != "" {
				thread, err := app.Client.CreateThread(context.Background(), openai.ThreadRequest{})
				if err != nil {
					c.Logger().Error(err)
					return internalServerError(c)
				}

				app.ThreadID = thread.ID
			}
		}

		return next(c)
	}
}

func (app *Application) ThreadValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.FileID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, jsonError(errors.New("file is not uploaded")))
		}

		if app.ThreadID == "" {
			thread, err := app.Client.CreateThread(context.Background(), openai.ThreadRequest{})
			if err != nil {
				c.Logger().Error(err)
				return internalServerError(c)
			}

			app.ThreadID = thread.ID
		}

		return next(c)
	}
}
