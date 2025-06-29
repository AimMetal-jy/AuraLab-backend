# 统一转录服务API接口文档

## 概述

本文档描述了统一的转录服务API接口，支持WhisperX和BlueLM两种模型。所有接口都使用统一的URL格式：`127.0.0.1:8888/model/`，通过查询参数区分不同的模型和操作。

## 基础URL格式

```
http://127.0.0.1:8888/model/?model={model_type}&action={action_type}&{additional_params}
```

### 参数说明

- `model`: 模型类型，支持 `whisperx` 或 `bluelm`
- `action`: 操作类型，支持 `submit`、`status`、`download`、`list`
- `additional_params`: 根据不同操作需要的额外参数

## API接口详情

### 1. 任务提交 (Submit)

#### WhisperX任务提交
```http
POST /model/?model=whisperx&action=submit
Content-Type: multipart/form-data

# Form Data:
# audio: 音频文件
# language: 语言代码 (可选)
# compute_type: 计算类型 (可选)
# batch_size: 批处理大小 (可选)
# chunk_size: 分块大小 (可选)
# return_char_alignments: 是否返回字符对齐 (可选)
```

#### BlueLM任务提交
```http
POST /model/?model=bluelm&action=submit
Content-Type: multipart/form-data

# Form Data:
# audio: 音频文件
```

#### 响应示例
```json
{
  "task_id": "task_20241201_123456_abc123",
  "status": "processing",
  "message": "Task submitted successfully"
}
```

### 2. 状态查询 (Status)

#### 查询任务状态
```http
GET /model/?model={whisperx|bluelm}&action=status&task_id={task_id}
```

#### 响应示例
```json
{
  "task_id": "task_20241201_123456_abc123",
  "status": "completed",
  "created_at": "2024-12-01T12:34:56Z",
  "completed_at": "2024-12-01T12:35:30Z",
  "result_file": "transcription_task_20241201_123456_abc123.txt"
}
```

#### 状态值说明
- `pending`: 等待处理
- `processing`: 正在处理
- `completed`: 处理完成
- `failed`: 处理失败

### 3. 文件下载 (Download)

#### WhisperX文件下载
```http
GET /model/?model=whisperx&action=download&task_id={task_id}&file_name={file_type}
```

**file_name参数值：**
- `transcription`: 转录文本文件
- `wordstamps`: 词级时间戳文件
- `diarization`: 说话人分离文件
- `speaker`: 说话人识别文件

#### BlueLM文件下载
```http
GET /model/?model=bluelm&action=download&task_id={task_id}
```

#### 响应
- 成功：返回文件内容（Content-Type根据文件类型设置）
- 失败：返回JSON错误信息

### 4. 任务列表 (List)

#### BlueLM任务列表
```http
GET /model/?model=bluelm&action=list&status={status}&limit={limit}
```

**可选参数：**
- `status`: 过滤特定状态的任务
- `limit`: 限制返回数量

#### 响应示例
```json
{
  "tasks": [
    {
      "task_id": "task_20241201_123456_abc123",
      "status": "completed",
      "created_at": "2024-12-01T12:34:56Z",
      "completed_at": "2024-12-01T12:35:30Z"
    }
  ],
  "total": 1
}
```

#### WhisperX任务列表
```http
GET /model/?model=whisperx&action=list
```

**注意：** WhisperX暂不支持任务列表功能，会返回501 Not Implemented。

## 完整使用示例

### 1. 提交BlueLM转录任务

```bash
curl -X POST "http://127.0.0.1:8888/model/?model=bluelm&action=submit" \
  -F "audio=@example.wav"
```

### 2. 查询任务状态

```bash
curl "http://127.0.0.1:8888/model/?model=bluelm&action=status&task_id=task_20241201_123456_abc123"
```

### 3. 下载转录结果

```bash
curl "http://127.0.0.1:8888/model/?model=bluelm&action=download&task_id=task_20241201_123456_abc123" \
  -o transcription_result.txt
```

### 4. 列出所有任务

```bash
curl "http://127.0.0.1:8888/model/?model=bluelm&action=list"
```

### 5. WhisperX示例

```bash
# 提交任务
curl -X POST "http://127.0.0.1:8888/model/?model=whisperx&action=submit" \
  -F "audio=@example.wav" \
  -F "language=zh"

# 查询状态
curl "http://127.0.0.1:8888/model/?model=whisperx&action=status&task_id=task_id_here"

# 下载转录文件
curl "http://127.0.0.1:8888/model/?model=whisperx&action=download&task_id=task_id_here&file_name=transcription" \
  -o whisperx_transcription.json
```

## 错误处理

### 常见错误响应

```json
{
  "error": "Bad Request",
  "message": "Model parameter is required (whisperx or bluelm)"
}
```

```json
{
  "error": "Not Found",
  "message": "Task not found"
}
```

```json
{
  "error": "Internal Server Error",
  "message": "Failed to process audio file"
}
```

## 向后兼容性

为了保持向后兼容，原有的API端点仍然可用：

- `POST /bluelm/transcription`
- `GET /bluelm/transcription/status/:task_id`
- `GET /bluelm/transcription/download/:task_id`
- `GET /bluelm/transcription/tasks`
- `POST /whisperx`
- `GET /whisperx/status/:task_id`
- `GET /whisperx/download/:task_id/:file_type`

## 注意事项

1. **文件大小限制**：音频文件大小限制为100MB
2. **支持格式**：支持常见音频格式（WAV、MP3、M4A等）
3. **任务过期**：任务结果文件保存7天后自动删除
4. **并发限制**：同时处理的任务数量有限制
5. **参数验证**：所有必需参数都会进行验证
6. **错误日志**：所有错误都会记录到服务器日志中

## 测试建议

1. 先测试健康检查接口：`GET /bluelm/health`
2. 使用小文件测试任务提交和状态查询
3. 验证文件下载功能
4. 测试错误场景（无效参数、不存在的任务等）
5. 验证向后兼容性