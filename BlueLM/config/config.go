package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体用于映射 config.yaml 的内容

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	VivoAI struct {
		AppID  string `yaml:"app_id"`
		AppKey string `yaml:"app_key"`
	} `yaml:"vivo_ai"`
	FilePaths struct {
		UploadDir   string `yaml:"upload_dir"`
		DownloadDir string `yaml:"download_dir"`
	} `yaml:"file_paths"`
	WhisperX struct {
		URL string `yaml:"url"`
	} `yaml:"whisperx"`
}

// LoadConfig 从指定的路径加载和解析YAML配置文件
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	// 优先使用环境变量
	if appID := os.Getenv("APPID"); appID != "" {
		config.VivoAI.AppID = appID
	}
	if appKey := os.Getenv("APPKEY"); appKey != "" {
		config.VivoAI.AppKey = appKey
	}

	// 验证必要的配置
	if config.VivoAI.AppID == "" || config.VivoAI.AppKey == "" {
		return nil, fmt.Errorf("APPID and APPKEY must be provided via environment variables or config file")
	}

	return config, nil
}