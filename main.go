package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/AimMetal-jy/AuraLab-backend/services"

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
			Vcn  string `json:"vcn,omitempty"`
		}
		requestBody.Mode = "TTS_MODE_HUMAN"
		requestBody.Vcn = "M24"
		requestBody.Text = "你好，这是蓝心大模型的音频生成功能。"
		// json传入
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		switch requestBody.Mode{
			case "TTS_MODE_SHORT":
				requestBody.Mode = "short_audio_synthesis_jovi"
			case "TTS_MODE_LONG":
				requestBody.Mode = "long_audio_synthesis_screen"
			case "TTS_MODE_HUMAN":
				requestBody.Mode = "tts_humanoid_lam"
			case "TTS_MODE_REPLICA":
				requestBody.Mode = "tts_replica" // 音色复刻专用
		}
		//调用蓝心大模型生成pcm切片
		res, e := app.TTS(requestBody.Mode, requestBody.Vcn, requestBody.Text)
		if e != nil {
			c.String(http.StatusInternalServerError, e.Error())
			fmt.Println("调用蓝心大模型生成pcm切片失败")
			return
		}
		outputFilePath := "output.wav"
		//将pcm切片转换为wav文件
		err := services.PcmtoWav(res, outputFilePath, 1, 16, 24000)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		//返回wav文件
		c.Header("Content-Type", "audio/wav")
		c.Header("Content-Disposition", "attachment; filename=output.wav")
		c.File(outputFilePath)
	})
	fmt.Println("服务启动成功")
	ginServer.Run(":8888")
}
