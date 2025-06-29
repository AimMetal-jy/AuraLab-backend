package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
	"github.com/dingdinglz/vivo"
	"github.com/gin-gonic/gin"
)

// TranscriptionHandler 处理长语音转写请求
func TranscriptionHandler(app *vivo.Vivo, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存上传的文件
		uploadFilePath, err := utils.SaveUploadedFile(c, cfg.FilePaths.UploadDir)
		if err != nil {
			utils.AbortWithBadRequest(c, err, "Failed to save uploaded file")
			return
		}
		//调用蓝心大模型长语音转写
		trans := app.NewTranscription(uploadFilePath)
		e := trans.Upload()
		if e != nil {
			utils.AbortWithInternalServerError(c, e)
			return
		}

		e = trans.Start()
		if e != nil {
			utils.AbortWithInternalServerError(c, e)
			return
		}

		taskID := getTaskID(trans)
		// 获取原始文件名
		filename := "unknown"
		if file, err := c.FormFile("file"); err == nil {
			filename = file.Filename
		}
		
		// 创建任务记录
		GlobalTaskManager.CreateTask(taskID, filename)
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusProcessing, "Processing started")

		go pollTranscriptionStatus(trans, cfg, taskID)

		c.JSON(http.StatusOK, gin.H{"task_id": taskID})
	}
}

// getTaskID uses reflection to access the private taskID field
func getTaskID(trans *vivo.Transcription) string {
	v := reflect.ValueOf(trans).Elem()
	taskIDField := v.FieldByName("taskID")
	if !taskIDField.IsValid() {
		return "unknown"
	}
	return taskIDField.String()
}

func pollTranscriptionStatus(trans *vivo.Transcription, cfg *config.Config, taskID string) {
	process := 0
	var e error
	for process != 100 {
		time.Sleep(1 * time.Second)
		// 查询任务进度
		process, e = trans.GetTaskInfo()
		if e != nil {
			utils.Log.Warnf("Failed to get task info for task %s: %v", taskID, e)
			GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusFailed, fmt.Sprintf("Error getting progress: %v", e))
			// 不中断轮询，继续尝试
			continue
		}
		utils.Log.Infof("Task %s progress: %d%%", taskID, process)
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusProcessing, fmt.Sprintf("Progress: %d%%", process))
	}

	result, e := trans.GetResult()
	if e != nil {
		utils.Log.Errorf("Failed to get result for task %s: %v", taskID, e)
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusFailed, fmt.Sprintf("Error getting result: %v", e))
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		utils.Log.Errorf("Failed to marshal result for task %s: %v", taskID, err)
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusFailed, fmt.Sprintf("Error serializing result: %v", err))
		return
	}

	filename := fmt.Sprintf("transcription_%s.json", taskID)
	downloadFilePath := filepath.Join(cfg.FilePaths.DownloadDir, filename)
	//将json数据写入文件
	err = os.WriteFile(downloadFilePath, jsonData, 0644)
	if err != nil {
		utils.Log.Errorf("Failed to write result to file for task %s: %v", taskID, err)
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusFailed, fmt.Sprintf("Error writing file: %v", err))
		return
	}

	// 更新任务状态为完成
	GlobalTaskManager.SetTaskFilePath(taskID, downloadFilePath)
	GlobalTaskManager.UpdateTaskStatus(taskID, TaskStatusCompleted, "Transcription completed successfully")
	utils.Log.Infof("Transcription task %s completed successfully. Result saved to %s", taskID, downloadFilePath)
}