package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "github.com/AimMetal-jy/AuraLab-backend/services"
	"github.com/dingdinglz/vivo"
	"github.com/gin-gonic/gin"
)

func main() {
	ginServer := gin.Default()
	app := vivo.NewVivoAIGC(vivo.Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})

	ginServer.POST("/bluelm/tts", func(c *gin.Context) {
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
			c.String(http.StatusBadRequest, err.Error())
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
		//调用蓝心大模型生成pcm切片
		res, e := app.TTS(requestBody.Mode, requestBody.Vcn, requestBody.Text)
		if e != nil {
			c.String(http.StatusInternalServerError, e.Error())
			fmt.Println("调用蓝心大模型生成pcm切片失败")
			return
		}
		fileName := time.Now().Format("20060102150405") + ".wav"
		downloadFilePath := "./file_io/download/" + "temp_" + fileName
		//将pcm切片转换为wav文件
		err := PcmtoWav(res, downloadFilePath, 1, 16, 24000)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		//返回wav文件
		c.Header("Content-Type", "audio/wav")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.File(downloadFilePath)
	})

	ginServer.POST("/bluelm/transcription", func(c *gin.Context) {
		//获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		uploadFileName := time.Now().Format("20060102150405") + ".wav"
		uploadFilePath := "./file_io/upload/" + "temp_" + uploadFileName
		err = c.SaveUploadedFile(file, uploadFilePath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		//调用蓝心大模型长语音转写
		trans := app.NewTranscription(uploadFilePath)
		e := trans.Upload()
		if e != nil {
			c.String(http.StatusInternalServerError, e.Error())
			return
		}

		e = trans.Start()
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		process := 0
		for process != 100 {
			time.Sleep(1 * time.Second)
			// 查询任务进度
			process, e = trans.GetTaskInfo()
			if e != nil {
				fmt.Println(e.Error())
				return
			}
			fmt.Println("当前任务进度：", process, "%")
		}
		result, e := trans.GetResult()
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		// for _, value := range result {
		// 	fmt.Println("开始秒数", value.Bg, "结束秒数", value.Ed, "内容：", value.Onebest)
		// }
		jsonData, err := json.Marshal(result)
		if err != nil {
			fmt.Println("JSON 编码错误:", err)
			return
		}
		downloadFileName := time.Now().Format("20060102150405") + ".json"
		downloadFilePath := "./file_io/download/" + "temp_" + downloadFileName
		//将json数据写入文件
		err = os.WriteFile(downloadFilePath, jsonData, 0644)
		if err != nil {
			fmt.Println("写入文件失败：", err)
			return
		}
		fmt.Println("JSON 数据已写入文件:", downloadFilePath)
		//返回json数据
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename="+downloadFileName)
		c.String(http.StatusOK, string(jsonData))
	})

	ginServer.POST("/whisperx", func(c *gin.Context) {
		// 获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "No file provided",
				"error":   err.Error(),
			})
			return
		}

		// 保存上传的文件到临时目录
		uploadFileName := time.Now().Format("20060102150405") + "_" + file.Filename
		uploadFilePath := "./file_io/upload/temp_" + uploadFileName
		err = c.SaveUploadedFile(file, uploadFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save uploaded file",
				"error":   err.Error(),
			})
			return
		}

		// 调用Flask WhisperX服务
		result, err := callWhisperXService(uploadFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "WhisperX processing failed",
				"error":   err.Error(),
			})
			return
		}

		// 清理临时文件
		defer func() {
			if _, err := os.Stat(uploadFilePath); err == nil {
				os.Remove(uploadFilePath)
			}
		}()

		// 返回结果
		c.JSON(http.StatusOK, result)
	})

	ginServer.POST("/bluelm/chat", func(ctx *gin.Context) {
		// 定义请求体结构
		var requestBody struct {
			Message   string                `json:"message"`   // 用户消息
			SessionID string                `json:"session_id"` // 会话ID（可选）
			History_Messages  []vivo.ChatMessage    `json:"history_messages"`   // 历史消息（可选）

		}

		// 解析JSON请求
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"message": err.Error(),
			})
			return
		}

		// 检查消息是否为空
		if requestBody.Message == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Message cannot be empty",
				"message": "请提供有效的消息内容",
			})
			return
		}

		// 生成或使用现有的会话ID
		session_id := requestBody.SessionID
		if session_id == "" {
			session_id = vivo.GenerateSessionID()
		}

		// 构建消息历史
		var history_messages []vivo.ChatMessage
		if len(requestBody.History_Messages) > 0 {
			// 使用提供的历史消息
			history_messages = requestBody.History_Messages
		}

		// 添加当前用户消息
		history_messages = append(history_messages, vivo.ChatMessage{
			Role:    vivo.CHAT_ROLE_USER,
			Content: requestBody.Message,
		})

		// 调用AI聊天接口
		res, err := app.Chat(vivo.GenerateRequestID(), session_id, history_messages, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Chat service error",
				"message": err.Error(),
			})
			return
		}

		// 将AI回复添加到消息历史
		history_messages = append(history_messages, res)

		// 返回响应
		ctx.JSON(http.StatusOK, gin.H{
			"success":    true,
			"message":    "Chat completed successfully",
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
			"session_id": session_id,
			"data": gin.H{
				"reply":    res.Content,
				"role":     res.Role,
				"messages": history_messages, // 返回完整的消息历史供前端维护状态
			},
		})
	})
	
	fmt.Println("服务启动成功")
	ginServer.Run(":8888")
}

// callWhisperXService 调用Flask WhisperX服务
func callWhisperXService(audioFilePath string) (map[string]interface{}, error) {
	// 创建multipart form数据
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件
	file, err := os.Open(audioFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// 发送POST请求到Flask服务
	req, err := http.NewRequest("POST", "http://localhost:5000/whisperx/process", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 300 * time.Second} // 5分钟超时
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Flask service returned error: %s", string(body))
	}

	// 解析响应JSON
	var uploadResult map[string]interface{}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %v", err)
	}

	// 检查上传是否成功
	if success, ok := uploadResult["success"].(bool); !ok || !success {
		return uploadResult, fmt.Errorf("upload failed: %v", uploadResult["message"])
	}

	// 获取任务ID
	taskID, ok := uploadResult["task_id"].(string)
	if !ok {
		return nil, fmt.Errorf("no task_id in response")
	}

	// 轮询任务状态直到完成
	for {
		time.Sleep(2 * time.Second) // 等待2秒再查询

		statusReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:5000/whisperx/status/%s", taskID), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create status request: %v", err)
		}

		statusResp, err := client.Do(statusReq)
		if err != nil {
			return nil, fmt.Errorf("failed to get task status: %v", err)
		}

		statusBody, err := io.ReadAll(statusResp.Body)
		statusResp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read status response: %v", err)
		}

		var statusResult map[string]interface{}
		err = json.Unmarshal(statusBody, &statusResult)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status response: %v", err)
		}

		status, ok := statusResult["status"].(string)
		if !ok {
			return nil, fmt.Errorf("no status in response")
		}

		switch status {
		case "completed":
			// 任务完成，返回结果
			if result, exists := statusResult["result"]; exists {
				return result.(map[string]interface{}), nil
			}
			return statusResult, nil
		case "failed":
			// 任务失败
			return statusResult, fmt.Errorf("processing failed: %v", statusResult["message"])
		case "processing", "queued":
			// 继续等待
			fmt.Printf("Task %s status: %s\n", taskID, status)
			continue
		default:
			return nil, fmt.Errorf("unknown task status: %s", status)
		}
	}
}
