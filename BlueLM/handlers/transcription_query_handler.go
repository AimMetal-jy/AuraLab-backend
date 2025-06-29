package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
)

// TranscriptionStatusHandler 查询转录任务状态
func TranscriptionStatusHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("task_id")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "Task ID is required")
			return
		}

		task, exists := GlobalTaskManager.GetTask(taskID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
				"task_id": taskID,
			})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

// TranscriptionDownloadHandler 下载转录结果文件
func TranscriptionDownloadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("task_id")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "Task ID is required")
			return
		}

		task, exists := GlobalTaskManager.GetTask(taskID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
				"task_id": taskID,
			})
			return
		}

		if task.Status != TaskStatusCompleted {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Task is not completed yet",
				"status": task.Status,
				"task_id": taskID,
			})
			return
		}

		// 构建文件路径
		filename := "transcription_" + taskID + ".json"
		filePath := filepath.Join(cfg.FilePaths.DownloadDir, filename)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Result file not found",
				"task_id": taskID,
			})
			return
		}

		// 设置响应头
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/json")

		// 发送文件
		c.File(filePath)
	}
}

// TranscriptionTasksHandler 列出所有转录任务
func TranscriptionTasksHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		status := c.Query("status")

		tasks := GlobalTaskManager.GetAllTasks()

		// 按状态过滤
		if status != "" {
			filteredTasks := make([]*TaskInfo, 0)
			for _, task := range tasks {
				if strings.EqualFold(string(task.Status), status) {
					filteredTasks = append(filteredTasks, task)
				}
			}
			tasks = filteredTasks
		}

		// 限制返回数量（简单实现）
		if len(tasks) > 50 { // 默认限制50个
			tasks = tasks[:50]
		}

		c.JSON(http.StatusOK, gin.H{
			"tasks": tasks,
			"total": len(tasks),
		})
	}
}