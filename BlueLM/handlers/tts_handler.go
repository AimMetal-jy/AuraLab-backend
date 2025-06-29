package handlers

import (
	"time"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
	"github.com/dingdinglz/vivo"
	"github.com/gin-gonic/gin"
)

// TTSHandler 处理文本到语音的转换请求
func TTSHandler(app *vivo.Vivo, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Mode string `json:"mode"`
			Text string `json:"text"`
			Vcn  string `json:"vcn"`
		}
		requestBody.Mode = "TTS_MODE_HUMAN"
		requestBody.Vcn = "M24"
		requestBody.Text = "你好，这是蓝心大模型的音频生成功能。"
		// json传入
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			utils.AbortWithBadRequest(c, err, "Invalid request body")
			return
		}
		switch requestBody.Mode {
		case "short":
			requestBody.Mode = "short_audio_synthesis_jovi"
		case "long":
			requestBody.Mode = "long_audio_synthesis_screen"
		case "human":
			requestBody.Mode = "tts_humanoid_lam"
		case "replica":
			requestBody.Mode = "tts_replica" // 音色复刻专用
		}
		// 检查配置是否为占位符
	if cfg.VivoAI.AppID == "YOUR_VIVO_APP_ID" || cfg.VivoAI.AppKey == "YOUR_VIVO_APP_KEY" {
		utils.Log.Errorf("TTS service configuration error: Vivo AI credentials not configured")
		c.JSON(500, gin.H{
			"message": "TTS service configuration error: Please configure valid Vivo AI credentials in config.yaml. See config.example.yaml for reference.",
			"error":   "Invalid or placeholder credentials detected",
			"details": "Current app_id and app_key are placeholder values. Please replace them with actual Vivo AI credentials.",
		})
		return
	}

	//调用蓝心大模型生成pcm切片
	res, e := app.TTS(requestBody.Mode, requestBody.Vcn, requestBody.Text)
	if e != nil {
		utils.Log.Errorf("TTS service error: %v", e)
		// 检查是否是配置问题
		if e.Error() == "invalid app_id or app_key" || e.Error() == "unauthorized" {
			c.JSON(500, gin.H{
				"message": "TTS service configuration error: Please check your Vivo AI credentials in config.yaml",
				"error": e.Error(),
				"details": "The provided Vivo AI credentials appear to be invalid. Please verify your app_id and app_key.",
			})
		} else if e.Error() == "websocket: bad handshake" {
			c.JSON(500, gin.H{
				"message": "TTS service connection error: Unable to establish WebSocket connection with Vivo AI service",
				"error": e.Error(),
				"details": "This usually indicates network connectivity issues or invalid credentials. Please check your internet connection and Vivo AI credentials.",
			})
		} else {
			c.JSON(500, gin.H{
				"message": "TTS service error: " + e.Error(),
				"error": e.Error(),
			})
		}
		return
	}
		fileName := time.Now().Format("20060102150405") + ".wav"
		downloadFilePath := cfg.FilePaths.DownloadDir + "temp_" + fileName
		//将pcm切片转换为wav文件
		err := utils.PcmtoWav(res, downloadFilePath, 1, 16, 24000)
		if err != nil {
			utils.AbortWithInternalServerError(c, err)
			return
		}
		//返回wav文件
		c.Header("Content-Type", "audio/wav")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.File(downloadFilePath)
	}
}