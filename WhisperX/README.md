# WhisperX Flask Service

这是一个基于Flask的WhisperX语音处理服务，提供语音转录、对齐和说话人分离功能。

## 功能特性

- 语音转录（Speech-to-Text）
- 词级时间戳对齐
- 说话人分离（Speaker Diarization）
- 异步处理支持
- RESTful API接口
- 文件上传和下载

## 安装依赖

```bash
pip install -r requirements.txt
```

## 环境变量

需要设置以下环境变量：

```bash
# HuggingFace Token (用于模型下载和说话人分离)
HF_WHISPERX=your_token_here
```

## 启动服务

**方式一：使用批处理脚本（推荐）**
```bash
# 从项目根目录运行
start_flask.bat
```

**方式二：直接启动**
```bash
# 使用指定的Python解释器
C:\ProgramData\anaconda3\envs\whisper\python.exe app.py
```

### 单独启动Flask服务
```bash
cd WhisperX
python app.py
```
服务将在 http://localhost:5000 启动

### 启动完整后端服务
运行根目录下的启动脚本：
```bash
start_services.bat
```
这将同时启动：
- WhisperX Flask服务（端口5000）
- BlueLM Go服务（端口8888）

## API接口

### 1. 健康检查
```
GET /health
```

### 2. 处理音频文件
```
POST /whisperx/process
Content-Type: multipart/form-data

参数：
- file: 音频文件（支持wav, mp3, mp4, avi, mov, flac, m4a）

响应：
{
  "success": true,
  "message": "File uploaded successfully, processing started",
  "task_id": "uuid",
  "filename": "example.wav"
}
```

### 3. 查询任务状态
```
GET /whisperx/status/{task_id}

响应：
{
  "success": true,
  "task_id": "uuid",
  "status": "completed",  // queued, processing, completed, failed
  "message": "Processing completed",
  "result": { ... }  // 仅在completed状态时包含
}
```

### 4. 获取处理结果
```
GET /whisperx/result/{task_id}

响应：
{
  "success": true,
  "message": "Audio processing completed successfully",
  "data": {
    "transcription": [...],  // 转录结果
    "language": "en",       // 检测到的语言
    "output_files": { ... }  // 输出文件路径
  }
}
```

### 5. 下载结果文件
```
GET /whisperx/download/{task_id}/{file_type}

file_type可选值：
- transcription: 基础转录结果
- wordstamps: 词级时间戳
- diarization: 说话人分离数据
- speaker: 带说话人标签的转录结果
```

### 6. 列出所有任务
```
GET /whisperx/tasks

响应：
{
  "success": true,
  "tasks": [...],
  "total": 10
}
```

## 通过Go服务调用

现在可以通过主Go服务（端口8888）调用WhisperX功能：

```
POST http://localhost:8888/whisperx
Content-Type: multipart/form-data

参数：
- file: 音频文件

响应：
{
  "success": true,
  "message": "Audio processing completed successfully",
  "data": {
    "transcription": [...],
    "language": "en"
  }
}
```

## 文件结构

```
WhisperX/
├── app.py                 # Flask应用主文件
├── whisperx_service.py    # WhisperX服务封装
├── whisperx_main.py       # 原始处理脚本
├── requirements.txt       # Python依赖
└── README.md             # 说明文档
```

## 注意事项

1. 确保有足够的GPU内存用于模型加载
2. 大文件处理可能需要较长时间
3. 说话人分离功能需要HuggingFace访问令牌
4. 临时文件会自动清理
5. 服务支持最大100MB文件上传

## 错误处理

服务包含完整的错误处理机制：
- 文件格式验证
- 文件大小限制
- 处理超时控制
- 自动清理临时文件
- 详细的错误信息返回