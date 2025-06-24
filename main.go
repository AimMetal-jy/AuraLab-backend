package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	fmt.Println("服务启动成功")
	ginServer.Run(":8888")
}
