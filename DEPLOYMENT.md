# AuraLab Backend 部署指南

## 概述

AuraLab Backend 现在包含两个主要服务：
- **BlueLM Go 服务** (端口 8888): 主要的 AI 服务，包括蓝心大模型的 TTS、ASR、Chat 功能
- **WhisperX Flask 服务** (端口 5000): 专门处理 WhisperX 语音转录和说话人分离

## 系统要求

### 软件依赖
- Go 1.19+
- Python 3.8+
- CUDA 支持的 GPU (推荐，用于 WhisperX)

### 环境变量
```bash
# HuggingFace Token (用于 WhisperX 模型下载)
HF_WHISPERX=your_token_here

# 蓝心大模型配置
APPID=your_app_id_here
APPKEY=your_app_key_here
```

## 安装步骤

### 1. 安装 Python 依赖
```bash
cd WhisperX
pip install -r requirements.txt
```

### 2. 安装 Go 依赖
```bash
cd BlueLM
go mod tidy
```

## 启动服务

### 方式一：使用批处理脚本（推荐）

**注意**: 启动脚本使用指定的Python解释器路径 `C:\ProgramData\anaconda3\envs\whisper\python.exe`，请确保该路径存在且环境已正确配置。

```bash
# 启动所有服务
start_services.bat

# 或者分别启动
start_flask.bat  # 启动 WhisperX Flask 服务
start_go.bat     # 启动 BlueLM Go 服务
```

### 方式二：手动启动

#### 启动 WhisperX Flask 服务
```bash
cd WhisperX
python app.py
```
服务将在 http://localhost:5000 启动

#### 启动 BlueLM Go 服务
```bash
cd BlueLM
go run main.go
```
服务将在 http://localhost:8888 启动

## 服务架构

```
用户请求 → Go 服务 (8888) → Flask 服务 (5000) → WhisperX 处理
                ↓
            统一响应格式
```

### 统一接口
所有 AI 功能都通过 Go 服务的 8888 端口提供：
- `/bluelm/tts` - 文本转语音
- `/bluelm/transcription` - 语音转文本
- `/bluelm/chat` - 对话
- `/whisperx` - WhisperX 语音处理（内部调用 Flask 服务）

### Flask 服务独立接口
WhisperX Flask 服务也提供独立的 RESTful API（端口 5000）：
- `/health` - 健康检查
- `/whisperx/process` - 处理音频文件
- `/whisperx/status/{task_id}` - 查询任务状态
- `/whisperx/result/{task_id}` - 获取处理结果
- `/whisperx/download/{task_id}/{file_type}` - 下载结果文件
- `/whisperx/tasks` - 列出所有任务

## 测试

### 运行集成测试
```bash
python test_integration.py
```

### 手动测试

#### 测试 WhisperX 功能
```bash
curl -X POST -F "file=@test_audio.wav" http://localhost:8888/whisperx
```

#### 测试健康检查
```bash
curl http://localhost:5000/health
curl http://localhost:8888/bluelm/chat -X POST -H "Content-Type: application/json" -d '{"message":"hello"}'
```

## 文件结构

```
AuraLab-backend/
├── BlueLM/                    # Go 主服务
│   ├── main.go               # 主服务入口
│   ├── services/
│   │   └── pcmtowav.go      # PCM转WAV工具
│   ├── file_io/
│   │   ├── upload/          # 文件上传目录
│   │   └── download/        # 文件下载目录
│   └── audio_examples/      # 音频示例文件
├── WhisperX/                  # Python Flask 服务
│   ├── app.py               # Flask 应用主文件
│   ├── whisperx_service.py  # WhisperX 服务封装
│   ├── whisperx_main.py     # 原始处理脚本
│   ├── requirements.txt     # Python 依赖
│   └── README.md           # 服务说明
├── start_services.bat        # 启动所有服务
├── start_flask.bat          # 启动 Flask 服务
├── start_go.bat            # 启动 Go 服务
├── test_integration.py      # 集成测试脚本
└── DEPLOYMENT.md           # 本文档
```

## 故障排除

### 常见问题

1. **Flask 服务启动失败**
   - 检查 Python 依赖是否安装完整
   - 确认端口 5000 未被占用
   - 检查 CUDA 环境（如果使用 GPU）

2. **Go 服务编译失败**
   - 运行 `go mod tidy` 更新依赖
   - 检查 Go 版本是否符合要求

3. **WhisperX 处理失败**
   - 确认 HuggingFace Token 设置正确
   - 检查音频文件格式是否支持
   - 查看 Flask 服务日志

4. **服务间通信失败**
   - 确认两个服务都已启动
   - 检查防火墙设置
   - 验证端口配置

### 日志查看
- Go 服务日志：控制台输出
- Flask 服务日志：控制台输出
- 详细错误信息：检查各服务的控制台输出

## 性能优化

1. **GPU 加速**：确保 CUDA 环境正确配置
2. **内存管理**：大文件处理时注意内存使用
3. **并发处理**：Flask 服务支持多任务并发
4. **文件清理**：定期清理临时文件

## 安全注意事项

1. 不要在生产环境中暴露 Flask 服务的 5000 端口
2. 设置适当的文件上传大小限制
3. 定期清理上传的临时文件
4. 保护 HuggingFace Token 等敏感信息