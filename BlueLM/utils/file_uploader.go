package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// SaveUploadedFile saves the uploaded file to the specified directory.
func SaveUploadedFile(c *gin.Context, uploadDir string) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	fileName := time.Now().Format("20060102150405") + "_" + file.Filename
	filePath := fmt.Sprintf("%s/%s", uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}