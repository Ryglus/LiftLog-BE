package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

// HandleFileUpload processes file uploads and returns the file path with a shortened UUID as filename
func HandleFileUpload(c *gin.Context, formKey string, dest string) (string, error) {
	file, header, err := c.Request.FormFile(formKey)
	if err != nil {
		return "", err // Handle missing file case in the caller
	}

	// Generate a shortened UUID for the image file
	fileExtension := filepath.Ext(header.Filename)
	var newFileName string
	var filePath string

	// Keep generating until a unique file name is found
	for {
		shortUUID := uuid.New().String()[:8] // Use only the first 8 characters of the UUID
		newFileName = fmt.Sprintf("%s%s", shortUUID, fileExtension)
		filePath = fmt.Sprintf("./uploads/%s/%s", dest, newFileName)

		// Check if the file already exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break // Exit the loop if the file doesn't exist
		}
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(fmt.Sprintf("./uploads/%s", dest), os.ModePerm); err != nil {
		return "", err
	}

	// Create file path for saving the file with shortened UUID as filename
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

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
