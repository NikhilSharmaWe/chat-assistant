package main

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sashabaranov/go-openai"
)

func (app *Application) Router() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.POST("upload", app.HandleUpload, app.DeleteOldFileIfExists)
	e.POST("chat", app.HandleMessage, app.ThreadValidator)

	return e
}

func (app *Application) HandleUpload(c echo.Context) error {
	uploadReq := UploadRequest{}
	if err := c.Bind(&uploadReq); err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	file, err := app.Client.CreateFile(context.Background(), openai.FileRequest{
		FileName: filepath.Base(uploadReq.FilePath),
		FilePath: uploadReq.FilePath,
		Purpose:  string(openai.PurposeAssistants),
	})
	if err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	_, err = app.Client.ModifyAssistant(context.Background(), app.AssistantID, openai.AssistantRequest{
		Name:         &assistantName,
		Description:  &description,
		Model:        openai.GPT4TurboPreview,
		Instructions: &instructions,
		Tools:        []openai.AssistantTool{{Type: openai.AssistantToolTypeRetrieval}},
		FileIDs:      []string{file.ID},
	})
	if err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	app.FileID = file.ID

	return c.JSON(http.StatusOK, map[string]string{"success": "file uploaded"})
}

func (app *Application) HandleMessage(c echo.Context) error {
	messageReq := MessageRequest{}
	if err := c.Bind(&messageReq); err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	_, err := app.Client.CreateMessage(context.Background(), app.ThreadID, openai.MessageRequest{
		Role:    "user",
		Content: messageReq.Message,
	})
	if err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	run, err := app.Client.CreateRun(context.Background(), app.ThreadID, openai.RunRequest{
		AssistantID: app.AssistantID,
	})
	if err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	var status string

	for {
		run, err := app.Client.RetrieveRun(context.Background(), app.ThreadID, run.ID)
		if err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		if run.Status == "completed" || run.Status == "failed" {
			status = string(run.Status)
			break
		}
	}

	switch status {
	case "completed":
		messages, err := app.Client.ListMessage(context.Background(), app.ThreadID, nil, nil, nil, nil)
		if err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		return c.JSON(http.StatusOK, map[string]string{"response": messages.Messages[0].Content[0].Text.Value})
	case "failed":
		return echo.NewHTTPError(http.StatusInternalServerError, jsonError(errors.New("failed to get the response")))
	}

	return nil
}
