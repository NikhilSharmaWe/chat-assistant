package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	assistantName = "File-Thinker"
	description   = "File-Thinker is a file master"
	instructions  = "You function as a file interpreter, providing responses to inquiries related to the uploaded file."
)

func internalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func jsonError(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
