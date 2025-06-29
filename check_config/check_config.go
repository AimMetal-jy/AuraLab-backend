package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration structure
type Config struct {
	VivoAI struct {
		AppID  string `yaml:"app_id"`
		AppKey string `yaml:"app_key"`
	} `yaml:"vivo_ai"`
}

func main() {
	fmt.Println("=== AuraLab Backend Configuration Checker ===")
	fmt.Println()

	// Check environment variables
	fmt.Println("1. Checking Environment Variables:")
	appID := os.Getenv("APPID")
	appKey := os.Getenv("APPKEY")
	hfToken := os.Getenv("HF_WHISPERX")

	if appID != "" && appKey != "" {
		fmt.Println("   ✓ APPID and APPKEY environment variables are set")
		fmt.Println("   → TTS service should work with environment variables")
	} else {
		fmt.Println("   ⚠ APPID and APPKEY environment variables not found")
		fmt.Println("   → Will check config.yaml file...")
	}

	if hfToken != "" {
		fmt.Println("   ✓ HF_WHISPERX environment variable is set")
	} else {
		fmt.Println("   ⚠ HF_WHISPERX environment variable not found")
		fmt.Println("   → WhisperX may have issues downloading models")
	}

	fmt.Println()

	// Check config.yaml file
	fmt.Println("2. Checking config.yaml file:")
	configPath := filepath.Join("..", "BlueLM", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("   ✗ config.yaml not found at %s\n", configPath)
		fmt.Println("   → Please copy config.example.yaml to config.yaml and configure it")
		return
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("   ✗ Error reading config.yaml: %v\n", err)
		return
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("   ✗ Error parsing config.yaml: %v\n", err)
		return
	}

	fmt.Println("   ✓ config.yaml file found and parsed successfully")

	// Check Vivo AI configuration
	if config.VivoAI.AppID == "YOUR_VIVO_APP_ID" || config.VivoAI.AppKey == "YOUR_VIVO_APP_KEY" {
		fmt.Println("   ⚠ Vivo AI credentials are still placeholder values")
		fmt.Println("   → TTS service will return 500 errors until configured")
		fmt.Println("   → Please replace YOUR_VIVO_APP_ID and YOUR_VIVO_APP_KEY with actual values")
	} else if config.VivoAI.AppID == "" || config.VivoAI.AppKey == "" {
		fmt.Println("   ⚠ Vivo AI credentials are empty")
		fmt.Println("   → TTS service will not work")
	} else {
		fmt.Println("   ✓ Vivo AI credentials appear to be configured")
		fmt.Println("   → TTS service should work (if credentials are valid)")
	}

	fmt.Println()

	// Summary
	fmt.Println("3. Configuration Summary:")
	if (appID != "" && appKey != "") || (config.VivoAI.AppID != "" && config.VivoAI.AppKey != "" && config.VivoAI.AppID != "YOUR_VIVO_APP_ID") {
		fmt.Println("   ✓ TTS Service: Ready")
	} else {
		fmt.Println("   ✗ TTS Service: Not configured")
	}

	if hfToken != "" {
		fmt.Println("   ✓ WhisperX Service: Ready")
	} else {
		fmt.Println("   ⚠ WhisperX Service: May have model download issues")
	}

	fmt.Println("   ✓ Chat Service: Ready (no additional config required)")
	fmt.Println("   ✓ BlueLM Transcription: Ready (no additional config required)")

	fmt.Println()
	fmt.Println("=== Configuration Check Complete ===")
	fmt.Println("For more information, see README.md and config.example.yaml")
}