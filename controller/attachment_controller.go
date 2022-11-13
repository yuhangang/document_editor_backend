package controller

import (
	constant "echoapp/commons"
	"echoapp/container"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

const (
	FileDirectory = "uploads/"
)

type AttachmentController interface {
	UploadAttachment(c echo.Context) error
	GetAttachment(c echo.Context) error
	DeleteAttachment(c echo.Context) error
}

// Products is a http.Handler
type attachmentController struct {
	l         *log.Logger
	container container.Container
}

func NewFileController(l *log.Logger, container container.Container) AttachmentController {
	return &attachmentController{l, container}
}

func (u *attachmentController) GetAttachment(c echo.Context) error {
	attachmentId := c.Param("id")
	return c.Attachment(attachmentId, path.Join(FileDirectory, attachmentId))
}

func (u *attachmentController) UploadAttachment(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("files")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path.Join(FileDirectory, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully.</p>", file.Filename))
}

func (u *attachmentController) DeleteAttachment(c echo.Context) error {
	attachmentId := c.Param("id")
	// Removing file from the directory
	// Using Remove() function
	e := os.Remove(path.Join(FileDirectory, attachmentId))
	if e != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s deleted successfully.</p>", attachmentId))
}
