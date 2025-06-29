# BlueLM Transcription API 测试指南

## 新增的 API 接口

### 1. 查询任务状态
```bash
GET /bluelm/transcription/status/{task_id}
```

**示例请求：**
```bash
curl -X GET "http://localhost:8888/bluelm/transcription/status/your_task_id"
```

**响应示例：**
```json
{
  "task_id": "your_task_id",
  "status": "completed",
  "message": "Transcription completed successfully",
  "created_at": "2024-01-15T10:30:00Z",
  "filename": "audio.wav"
}
```

**状态说明：**
- `pending`: 任务已创建，等待开始
- `processing`: 任务正在处理中
- `completed`: 任务已完成
- `failed`: 任务失败

### 2. 下载转录结果
```bash
GET /bluelm/transcription/download/{task_id}
```

**示例请求：**
```bash
curl -X GET "http://localhost:8888/bluelm/transcription/download/your_task_id" -o transcription_result.json
```

**注意：**
- 只有状态为 `completed` 的任务才能下载
- 返回 JSON 格式的转录结果文件

### 3. 列出所有任务
```bash
GET /bluelm/transcription/tasks
```

**示例请求：**
```bash
# 获取所有任务
curl -X GET "http://localhost:8888/bluelm/transcription/tasks"

# 按状态过滤
curl -X GET "http://localhost:8888/bluelm/transcription/tasks?status=completed"
```

**响应示例：**
```json
{
  "tasks": [
    {
      "task_id": "task_1",
      "status": "completed",
      "message": "Transcription completed successfully",
      "created_at": "2024-01-15T10:30:00Z",
      "filename": "audio1.wav"
    },
    {
      "task_id": "task_2",
      "status": "processing",
      "message": "Progress: 75%",
      "created_at": "2024-01-15T11:00:00Z",
      "filename": "audio2.mp3"
    }
  ],
  "total": 2
}
```

## 完整的使用流程

### 1. 提交转录任务
```bash
curl -X POST "http://localhost:8888/bluelm/transcription" \
  -F "file=@/path/to/your/audio.wav"
```

**响应：**
```json
{
  "task_id": "generated_task_id"
}
```

### 2. 查询任务状态
```bash
curl -X GET "http://localhost:8888/bluelm/transcription/status/generated_task_id"
```

### 3. 下载结果（任务完成后）
```bash
curl -X GET "http://localhost:8888/bluelm/transcription/download/generated_task_id" \
  -o transcription_result.json
```

## 错误处理

### 任务不存在
```json
{
  "error": "Task not found",
  "task_id": "invalid_task_id"
}
```

### 任务未完成
```json
{
  "error": "Task is not completed yet",
  "status": "processing",
  "task_id": "your_task_id"
}
```

### 文件不存在
```json
{
  "error": "Result file not found",
  "task_id": "your_task_id"
}
```

## 注意事项

1. **文件命名规则**：结果文件现在使用 `transcription_{task_id}.json` 格式命名
2. **任务状态管理**：所有任务状态都在内存中管理，服务重启后会丢失
3. **文件存储**：转录结果保存在配置文件指定的 `download_dir` 目录中
4. **并发安全**：任务管理器使用读写锁确保并发安全
5. **任务限制**：列出任务接口默认最多返回50个任务

## 服务启动

确保服务在端口 8888 上运行：
```bash
cd BlueLM
go run main.go
```

服务启动后会显示：
```
服务启动成功，监听端口 :8888
```