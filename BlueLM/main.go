package main

import (
	"fmt"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/config"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/handlers"
	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
	"github.com/dingdinglz/vivo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		utils.Log.Fatalf("Failed to load config: %v", err)
	}

	ginServer := gin.Default()

	// 配置CORS中间件
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	ginServer.Use(cors.New(corsConfig))
	app := vivo.NewVivoAIGC(vivo.Config{
		AppID:  cfg.VivoAI.AppID,
		AppKey: cfg.VivoAI.AppKey,
	})

	// Register handlers
	// Health check endpoint
	ginServer.GET("/bluelm/health", handlers.HealthHandler)
	
	// Legacy endpoints (保持向后兼容)
	ginServer.POST("/bluelm/tts", handlers.TTSHandler(app, cfg))
	ginServer.POST("/bluelm/transcription", handlers.TranscriptionHandler(app, cfg))
	ginServer.POST("/bluelm/chat", handlers.ChatHandler(app))
	ginServer.POST("/whisperx", handlers.WhisperXHandler(cfg))
	ginServer.GET("/bluelm/transcription/status/:task_id", handlers.TranscriptionStatusHandler(cfg))
	ginServer.GET("/bluelm/transcription/download/:task_id", handlers.TranscriptionDownloadHandler(cfg))
	ginServer.GET("/bluelm/transcription/tasks", handlers.TranscriptionTasksHandler(cfg))
	// WhisperX状态和下载接口现在通过统一API提供
	
	// 统一的模型API接口
	ginServer.POST("/model", handlers.UnifiedModelHandler(app, cfg))
	ginServer.GET("/model", handlers.UnifiedModelHandler(app, cfg))

	fmt.Println("服务启动成功，监听端口", cfg.Server.Port)
	ginServer.Run(cfg.Server.Port)
}
