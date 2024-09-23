package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

// HandleFileUpload processes file uploads and returns the file path with a UUID as filename
func HandleFileUpload(c *gin.Context, formKey string, dest string) (string, error) {
	file, header, err := c.Request.FormFile(formKey)
	if err != nil {
		return "", err // Handle missing file case in the caller
	}

	// Generate a new UUID for the image file
	fileExtension := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return "", err
	}

	// Create file path for saving the file with UUID as filename
	filePath := fmt.Sprintf("./uploads/%s/%s", dest, newFileName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the file to the uploads folder
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DeleteOldFile deletes the old image file
func DeleteOldFile(filePath string) error {
	if filePath != "" {
		err := os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
