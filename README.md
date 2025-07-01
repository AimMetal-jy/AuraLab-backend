# AuraLab-backend

欢迎使用 AuraLab-backend，这是一个集成了 Go 和 Python 微服务的后端系统，旨在提供强大的 AI 功能，包括语音处理和大型语言模型交互。

## 概述

AuraLab Backend 包含两个核心服务：
- **BlueLM Go 服务** (端口 `8888`): 作为主服务网关，负责处理用户请求、调用蓝心大模型功能，并与 WhisperX 服务通信。
- **WhisperX Flask 服务** (端口 `5000`): 专门用于高性能的语音转录、对齐和说话人分离。

## 主要功能

- [x] **文本转语音 (TTS)**: 将文本转换为自然的语音。
- [x] **长语音转写 - BlueLM**: 使用蓝心大模型进行语音转文本。
- [x] **长语音转写 - WhisperX**: 使用 WhisperX 进行高精度语音转录、对齐和说话人分离。
- [x] **AI 对话**: 与先进的语言模型进行交互。

## 系统架构

所有外部请求都通过 Go 服务进行路由，该服务根据需要将特定任务（如 WhisperX 处理）转发到 Python 服务。

```
用户请求 → Go 服务 (8888) → Flask 服务 (5000) → WhisperX 处理
                ↓
            统一响应格式
```

## 系统要求

- **Go**: 1.19+
- **Python**: 3.8+
- **CUDA**: 推荐使用支持 CUDA 的 GPU 以获得最佳 WhisperX 性能。

## 环境变量

在运行服务之前，请确保设置了以下环境变量：

```bash
# HuggingFace Token (用于 WhisperX 模型下载)
HF_WHISPERX=your_huggingface_token_here

# 蓝心大模型配置（优先级：环境变量 > config.yaml）
APPID=your_vivo_app_id_here
APPKEY=your_vivo_app_key_here
```

**配置优先级说明：**
- 系统会优先使用环境变量 `APPID` 和 `APPKEY`
- 如果环境变量不存在，则回退到 `config.yaml` 文件中的配置
- 如果两者都没有配置，服务将启动失败并报错

## 安装步骤

1.  **安装 Python 依赖**

    ```bash
    cd WhisperX
    pip install -r requirements.txt
    ```

2.  **安装 Go 依赖**

    ```bash
    cd BlueLM
    go mod tidy
    ```

## 启动服务

推荐使用根目录下的批处理脚本一键启动所有服务。

```bash
# 启动所有服务 (推荐)
start_services.bat
```

这将同时在新的终端窗口中启动 Go 服务和 Flask 服务。

你也可以单独启动每个服务：

```bash
# 单独启动 WhisperX Flask 服务
start_flask.bat

# 单独启动 BlueLM Go 服务
start_go.bat
```

## API 端点

### 统一接口 (Go 服务 @ 8888)

#### 1. 文本转语音 (TTS)
**接口：** `POST /bluelm/tts`

**请求示例：**
```bash
curl -X POST http://localhost:8888/bluelm/tts \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "human",
    "text": "你好，这是蓝心大模型的音频生成功能。",
    "vcn": "M24"
  }' \
  --output output.wav
```

**请求参数：**
- `mode`: 合成模式 (short/long/human/replica)
- `text`: 要转换的文本
- `vcn`: 音色选择 (M24等)

**响应：** 直接返回WAV音频文件流

#### 2. 语音转文本 (ASR)
**接口：** `POST /bluelm/transcription`

**请求示例：**
```bash
curl -X POST http://localhost:8888/bluelm/transcription \
  -F "file=@audio.wav"
```

**响应：**
```json
{
  "task_id": "generated_task_id_string"
}
```

**说明：** 异步处理，结果自动保存到下载目录

#### 3. AI对话
**接口：** `POST /bluelm/chat`

**请求示例：**
```bash
curl -X POST http://localhost:8888/bluelm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下自己",
    "session_id": "optional_session_id",
    "history_messages": []
  }'
```

**响应：**
```json
{
  "success": true,
  "message": "Chat completed successfully",
  "timestamp": "2024-01-01 12:00:00",
  "session_id": "generated_session_id",
  "data": {
    "reply": "AI回复内容",
    "role": "assistant",
    "messages": []
  }
}
```

#### 4. WhisperX语音处理 (代理接口)
**接口：** `POST /whisperx`

**请求示例：**
```bash
curl -X POST http://localhost:8888/whisperx \
  -F "file=@audio.wav"
```

**响应：**
```json
{
  "task_id": "uuid_generated_task_id"
}
```

**说明：** 代理到Flask服务进行处理，支持多种音频格式

### Flask 服务独立接口 (@ 5000)

#### 1. 健康检查
**接口：** `GET /health`

**请求示例：**
```bash
curl http://localhost:5000/health
```

**响应：**
```json
{
  "status": "ok",
  "message": "WhisperX service is running",
  "timestamp": "2024-01-01 12:00:00"
}
```

#### 2. 处理音频文件
**接口：** `POST /whisperx/process`

**请求示例：**
```bash
curl -X POST http://localhost:5000/whisperx/process \
  -F "file=@audio.wav"
```

**响应：**
```json
{
  "success": true,
  "message": "File uploaded successfully, processing started",
  "task_id": "uuid_task_id",
  "filename": "audio.wav"
}
```

#### 3. 查询任务状态
**接口：** `GET /whisperx/status/{task_id}`

**请求示例：**
```bash
curl http://localhost:5000/whisperx/status/your_task_id
```

**响应：**
```json
{
  "success": true,
  "task_id": "uuid_task_id",
  "status": "completed",
  "message": "Processing completed successfully",
  "created_at": 1640995200.0,
  "filename": "audio.wav",
  "result": {
    "success": true,
    "output_files": [
      "whisperx_output.json",
      "wordstamps.json",
      "diarization.json",
      "assign_speaker.json"
    ]
  }
}
```

#### 4. 下载结果文件
**接口：** `GET /whisperx/download/{task_id}/{file_type}`

**文件类型：**
- `transcription`: 基础转录结果
- `wordstamps`: 词级时间戳
- `diarization`: 说话人分离
- `speaker`: 说话人分配

**请求示例：**
```bash
curl http://localhost:5000/whisperx/download/your_task_id/transcription \
  --output result.json
```

#### 5. 列出所有任务
**接口：** `GET /whisperx/tasks`

**请求示例：**
```bash
curl http://localhost:5000/whisperx/tasks
```

**响应：**
```json
{
  "success": true,
  "tasks": [
    {
      "task_id": "uuid_task_id",
      "status": "completed",
      "message": "Processing completed successfully",
      "created_at": 1640995200.0,
      "filename": "audio.wav"
    }
  ],
  "total": 1
}
```

## 文件结构

```
AuraLab-backend/
├── BlueLM/                    # Go 主服务
│   ├── main.go               # 主服务入口
│   └── services/
│       └── pcmtowav.go      # PCM转WAV工具
├── WhisperX/                  # Python Flask 服务
│   ├── app.py               # Flask 应用主文件
│   ├── whisperx_service.py  # WhisperX 服务封装
│   └── requirements.txt     # Python 依赖
├── file_io/
│   ├── upload/              # 统一文件上传目录
│   └── download/            # 统一文件下载目录
├── .gitignore                 # Git 忽略文件配置
├── README.md                  # 本文档
├── start_services.bat        # 启动所有服务
├── start_flask.bat          # 启动 Flask 服务
├── start_go.bat            # 启动 Go 服务
└── README.md               # 本文档
```

## 测试

### 集成测试

运行所有服务后，可以通过以下方式测试：

#### 1. 健康检查测试
```bash
# 测试Flask服务健康状态
curl http://localhost:5000/health

# 预期响应：
# {
#   "status": "ok",
#   "message": "WhisperX service is running",
#   "timestamp": "2024-01-01 12:00:00"
# }
```

#### 2. TTS功能测试
```bash
# 测试文本转语音
curl -X POST http://localhost:8888/bluelm/tts \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "human",
    "text": "这是一个测试音频",
    "vcn": "M24"
  }' \
  --output test_output.wav

# 检查生成的音频文件
ls -la test_output.wav
```

#### 3. ASR功能测试
```bash
# 测试语音转文本（需要准备音频文件）
curl -X POST http://localhost:8888/bluelm/transcription \
  -F "file=@test_audio.wav"

# 预期响应：
# {
#   "task_id": "generated_task_id_string"
# }

# 检查下载目录中的结果文件
ls -la file_io/downloads/
```

#### 4. AI对话测试
```bash
# 测试AI聊天功能
curl -X POST http://localhost:8888/bluelm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下自己",
    "session_id": "test_session_001",
    "history_messages": []
  }'

# 预期响应：
# {
#   "success": true,
#   "message": "Chat completed successfully",
#   "timestamp": "2024-01-01 12:00:00",
#   "session_id": "test_session_001",
#   "data": {
#     "reply": "你好！我是蓝心大模型...",
#     "role": "assistant",
#     "messages": [...]
#   }
# }
```

#### 5. WhisperX功能测试
```bash
# 通过Go代理接口测试
curl -X POST http://localhost:8888/whisperx \
  -F "file=@test_audio.wav"

# 预期响应：
# {
#   "task_id": "uuid_generated_task_id"
# }

# 直接测试Flask服务
curl -X POST http://localhost:5000/whisperx/process \
  -F "file=@test_audio.wav"

# 查询任务状态
curl http://localhost:5000/whisperx/status/your_task_id

# 下载结果文件
curl http://localhost:5000/whisperx/download/your_task_id/transcription \
  --output whisperx_result.json
```

#### 6. 完整工作流测试
```bash
#!/bin/bash
# 完整测试脚本示例

echo "=== 开始集成测试 ==="

# 1. 健康检查
echo "1. 测试健康检查..."
curl -s http://localhost:5000/health | jq .

# 2. TTS测试
echo "2. 测试TTS功能..."
curl -X POST http://localhost:8888/bluelm/tts \
  -H "Content-Type: application/json" \
  -d '{"mode":"short","text":"测试音频","vcn":"M24"}' \
  --output test.wav

if [ -f "test.wav" ]; then
  echo "TTS测试成功，音频文件已生成"
else
  echo "TTS测试失败"
fi

# 3. 使用生成的音频测试ASR
echo "3. 测试ASR功能..."
ASR_RESULT=$(curl -X POST http://localhost:8888/bluelm/transcription \
  -F "file=@test.wav" -s)
echo "ASR结果: $ASR_RESULT"

# 4. AI对话测试
echo "4. 测试AI对话..."
curl -X POST http://localhost:8888/bluelm/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"你好"}' -s | jq .

echo "=== 测试完成 ==="
```

### 手动测试

#### 使用Postman测试
1. 导入API集合（可创建Postman Collection）
2. 设置环境变量：
   - `base_url_go`: http://localhost:8888
   - `base_url_flask`: http://localhost:5000
3. 按顺序测试各个接口

#### 使用curl脚本测试
创建测试脚本 `test_api.sh`：
```bash
#!/bin/bash
# API测试脚本

BASE_URL_GO="http://localhost:8888"
BASE_URL_FLASK="http://localhost:5000"

# 测试函数
test_endpoint() {
    echo "Testing: $1"
    echo "Command: $2"
    eval $2
    echo "---"
}

# 执行测试
test_endpoint "Health Check" "curl -s $BASE_URL_FLASK/health | jq ."
test_endpoint "AI Chat" "curl -X POST $BASE_URL_GO/bluelm/chat -H 'Content-Type: application/json' -d '{\"message\":\"Hello\"}' -s | jq ."
# 添加更多测试...
```

### 性能测试

#### 并发测试
```bash
# 使用ab工具进行并发测试
ab -n 100 -c 10 http://localhost:5000/health

# 使用wrk进行压力测试
wrk -t12 -c400 -d30s http://localhost:5000/health
```

#### 内存和CPU监控
```bash
# 监控Go服务
top -p $(pgrep -f "go run main.go")

# 监控Python服务
top -p $(pgrep -f "python app.py")

# 查看端口占用
netstat -tulpn | grep -E ":(5000|8888)"
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