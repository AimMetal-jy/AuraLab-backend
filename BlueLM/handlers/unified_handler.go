package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/dingdinglz/vivo"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
)

// UnifiedModelHandler 统一的模型处理入口
// 支持的URL格式：
// POST /model/?model=whisperx&action=submit
// POST /model/?model=bluelm&action=submit
// GET /model/?model=whisperx&action=status&task_id=xxx
// GET /model/?model=bluelm&action=status&task_id=xxx
// GET /model/?model=whisperx&action=download&task_id=xxx&file_name=xxx
// GET /model/?model=bluelm&action=download&task_id=xxx
// GET /model/?model=whisperx&action=list
// GET /model/?model=bluelm&action=list
func UnifiedModelHandler(app *vivo.Vivo, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取模型类型
		model := strings.ToLower(c.Query("model"))
		if model == "" {
			utils.AbortWithBadRequest(c, nil, "Model parameter is required (whisperx or bluelm)")
			return
		}

		// 获取操作类型
		action := strings.ToLower(c.Query("action"))
		if action == "" {
			utils.AbortWithBadRequest(c, nil, "Action parameter is required (submit, status, download, list)")
			return
		}

		// 根据模型类型分发请求
		switch model {
		case "whisperx":
			handleWhisperXRequest(c, cfg, action)
		case "bluelm":
			handleBlueLMRequest(c, app, cfg, action)
		default:
			utils.AbortWithBadRequest(c, nil, "Unsupported model. Supported models: whisperx, bluelm")
		}
	}
}

// handleWhisperXRequest 处理WhisperX相关请求
func handleWhisperXRequest(c *gin.Context, cfg *config.Config, action string) {
	switch action {
	case "submit":
		// 调用原有的WhisperX提交处理器
		WhisperXHandler(cfg)(c)
	case "status":
		// 处理状态查询
		taskID := c.Query("task_id")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "task_id parameter is required")
			return
		}
		handleWhisperXStatus(c, cfg, taskID)
	case "download":
		// 处理文件下载
		taskID := c.Query("task_id")
		fileName := c.Query("file_name")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "task_id parameter is required")
			return
		}
		if fileName == "" {
			utils.AbortWithBadRequest(c, nil, "file_name parameter is required (transcription, wordstamps, diarization, speaker)")
			return
		}
		handleWhisperXDownload(c, cfg, taskID, fileName)
	case "list":
		// 处理任务列表查询
		handleWhisperXList(c, cfg)
	default:
		utils.AbortWithBadRequest(c, nil, "Unsupported action for WhisperX. Supported actions: submit, status, download, list")
	}
}

// handleBlueLMRequest 处理BlueLM相关请求
func handleBlueLMRequest(c *gin.Context, app *vivo.Vivo, cfg *config.Config, action string) {
	switch action {
	case "submit":
		// 调用原有的BlueLM转录处理器
		TranscriptionHandler(app, cfg)(c)
	case "status":
		// 处理状态查询
		taskID := c.Query("task_id")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "task_id parameter is required")
			return
		}
		handleBlueLMStatus(c, cfg, taskID)
	case "download":
		// 处理文件下载
		taskID := c.Query("task_id")
		if taskID == "" {
			utils.AbortWithBadRequest(c, nil, "task_id parameter is required")
			return
		}
		handleBlueLMDownload(c, cfg, taskID)
	case "list":
		// 处理任务列表查询
		handleBlueLMList(c, cfg)
	default:
		utils.AbortWithBadRequest(c, nil, "Unsupported action for BlueLM. Supported actions: submit, status, download, list")
	}
}

// WhisperX相关的具体处理函数
func handleWhisperXStatus(c *gin.Context, cfg *config.Config, taskID string) {
	// 直接调用WhisperX服务的状态查询API
	statusURL := fmt.Sprintf("%s/whisperx/status/%s", cfg.WhisperX.URL, taskID)
	resp, err := http.Get(statusURL)
	if err != nil {
		utils.AbortWithInternalServerError(c, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.AbortWithInternalServerError(c, err)
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

func handleWhisperXDownload(c *gin.Context, cfg *config.Config, taskID, fileName string) {
	// 直接调用WhisperX服务的下载API
	downloadURL := fmt.Sprintf("%s/whisperx/download/%s/%s", cfg.WhisperX.URL, taskID, fileName)
	resp, err := http.Get(downloadURL)
	if err != nil {
		utils.AbortWithInternalServerError(c, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.AbortWithInternalServerError(c, err)
		return
	}

	// 设置适当的Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Data(resp.StatusCode, contentType, body)
}

func handleWhisperXList(c *gin.Context, _ *config.Config) {
	// WhisperX暂时没有列表功能，返回提示信息
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "List functionality not implemented for WhisperX",
		"message": "WhisperX service does not support task listing",
	})
}

// BlueLM相关的具体处理函数
func handleBlueLMStatus(c *gin.Context, cfg *config.Config, taskID string) {
	// 构造原有的路径参数格式
	c.Params = gin.Params{{Key: "task_id", Value: taskID}}
	// 调用原有的状态查询处理器
	TranscriptionStatusHandler(cfg)(c)
}

func handleBlueLMDownload(c *gin.Context, cfg *config.Config, taskID string) {
	// 构造原有的路径参数格式
	c.Params = gin.Params{{Key: "task_id", Value: taskID}}
	// 调用原有的下载处理器
	TranscriptionDownloadHandler(cfg)(c)
}

func handleBlueLMList(c *gin.Context, cfg *config.Config) {
	// 调用原有的任务列表处理器
	TranscriptionTasksHandler(cfg)(c)
}