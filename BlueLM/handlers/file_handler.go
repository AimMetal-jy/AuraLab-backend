package handlers

import (
	"net/http";
	"os";

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config";
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils";
	"github.com/gin-gonic/gin";
)

// FileInfo represents basic information about a file.
// We can extend this later to include metadata from a database.
type FileInfo struct {
	Name string `json:"name"`;
	Size int64  `json:"size"`;
}

// FilesHandler handles the request to list files.
func FilesHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadDir := cfg.FilePaths.UploadDir
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			utils.Log.Errorf("Failed to read directory %s: %v", uploadDir, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list files"})
			return
		}

		var fileInfos []FileInfo
		for _, file := range files {
			if !file.IsDir() {
				info, err := file.Info()
				if (err == nil) {
					fileInfos = append(fileInfos, FileInfo{Name: file.Name(), Size: info.Size()})
				}
			}
		}

		c.JSON(http.StatusOK, fileInfos)
	}
}