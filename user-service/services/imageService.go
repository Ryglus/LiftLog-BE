package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// HandleFileUpload processes file uploads and returns the file path
func HandleFileUpload(c *gin.Context, formKey string) (string, error) {
	file, header, err := c.Request.FormFile(formKey)
	if err != nil {
		return "", err // Handle missing file case in the caller
	}

	// Create file path for saving the file
	filePath := fmt.Sprintf("./uploads/%s", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy file to the uploads folder
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
