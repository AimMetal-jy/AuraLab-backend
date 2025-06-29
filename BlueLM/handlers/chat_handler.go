package handlers

import (
	"net/http"
	"time"

	"github.com/AimMetal-jy/AuraLab-backend/BlueLM/utils"
	"github.com/dingdinglz/vivo"
	"github.com/gin-gonic/gin"
)

// ChatHandler handles the AI chat requests.
func ChatHandler(app *vivo.Vivo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 定义请求体结构
		var requestBody struct {
			Message         string             `json:"message"`
			SessionID       string             `json:"session_id,omitempty"`
			HistoryMessages []vivo.ChatMessage `json:"history_messages,omitempty"`
		}

		// 解析JSON请求
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			utils.AbortWithBadRequest(ctx, err, "Invalid request format")
			return
		}

		// 检查消息是否为空
		if requestBody.Message == "" {
			utils.AbortWithBadRequest(ctx, nil, "Message cannot be empty")
			return
		}

		// 生成或使用现有的会话ID
		sessionID := requestBody.SessionID
		if sessionID == "" {
			sessionID = vivo.GenerateSessionID()
		}

		// 构建消息历史
		var historyMessages []vivo.ChatMessage
		if len(requestBody.HistoryMessages) > 0 {
			// 使用提供的历史消息
			historyMessages = requestBody.HistoryMessages
		}

		// 添加当前用户消息
		historyMessages = append(historyMessages, vivo.ChatMessage{
			Role:    vivo.CHAT_ROLE_USER,
			Content: requestBody.Message,
		})

		// 调用蓝心大模型
		res, err := app.Chat(vivo.GenerateSessionID(), sessionID, historyMessages, nil)
		if err != nil {
			utils.AbortWithInternalServerError(ctx, err)
			return
		}

		// 将AI回复添加到消息历史
		historyMessages = append(historyMessages, res)

		// 返回响应
		ctx.JSON(http.StatusOK, gin.H{
			"success":    true,
			"message":    "Chat completed successfully",
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
			"session_id": sessionID,
			"data": gin.H{
				"reply":    res.Content,
				"role":     res.Role,
				"messages": historyMessages, // 返回完整的消息历史供前端维护状态
			},
		})
	}
}