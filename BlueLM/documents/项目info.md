## 后端路由框架

URL="127.0.0.1"

Port="8888"

```
URL:Port
	|--/bluelm/transcription(蓝心大模型-长语音转写)
	|--/bluelm/chat(文本生成)
	|--/bluelm/tts(音频生成)
	|--/whisperx(whisperx-长语音转写)
```

## 接口详细信息

### 1. /bluelm/tts (音频生成)

**请求方法**: POST  
**Content-Type**: application/json

**请求参数**:
```json
{
  "mode": "human",     // 模式: "short", "long", "human", "replica"
  "text": "你好，这是蓝心大模型的音频生成功能。",  // 要转换的文本
  "vcn": "M24"         // 音色选择
}
```

**响应**: 返回生成的WAV音频文件

**使用示例**:
```bash
# curl示例
curl -X POST http://127.0.0.1:8888/bluelm/tts \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "human",
    "text": "你好，欢迎使用语音合成服务",
    "vcn": "M24"
  }' \
  --output output.wav
```

```javascript
// JavaScript示例
fetch('http://127.0.0.1:8888/bluelm/tts', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    mode: 'human',
    text: '你好，欢迎使用语音合成服务',
    vcn: 'M24'
  })
})
.then(response => response.blob())
.then(blob => {
  const url = URL.createObjectURL(blob);
  const audio = new Audio(url);
  audio.play();
});
```

### 2. /bluelm/transcription (蓝心大模型-长语音转写)

**请求方法**: POST  
**Content-Type**: multipart/form-data

**请求参数**:
- `file`: 音频文件 (支持WAV格式)

**响应**: 返回转写结果的JSON数据
```json
[
  {
    "bg": 0.0,        // 开始时间(秒)
    "ed": 2.5,        // 结束时间(秒)
    "onebest": "你好"  // 转写文本
  }
]
```

**使用示例**:
```bash
# curl示例
curl -X POST http://127.0.0.1:8888/bluelm/transcription \
  -F "file=@audio.wav"
```

```javascript
// JavaScript示例
const formData = new FormData();
formData.append('file', audioFile);

fetch('http://127.0.0.1:8888/bluelm/transcription', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('转写结果:', data);
});
```

### 3. /whisperx (WhisperX-长语音转写)

**请求方法**: POST  
**Content-Type**: multipart/form-data

**请求参数**:
- `file`: 音频文件 (推荐WAV格式)

**响应**: 返回详细的转写和说话人识别结果
```json
{
  "message": "音频处理完成",
  "timestamp": "20240101120000",
  "data": {
    "whisperx_output": {},    // 基础转写结果
    "wordstamps": {},         // 词级时间戳
    "assign_speaker": {},     // 说话人分配结果
    "diarization": {}         // 说话人分离结果
  }
}
```

**使用示例**:
```bash
# curl示例
curl -X POST http://127.0.0.1:8888/whisperx \
  -F "file=@audio.wav"
```

```javascript
// JavaScript示例
const formData = new FormData();
formData.append('file', audioFile);

fetch('http://127.0.0.1:8888/whisperx', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('WhisperX处理结果:', data);
});
```

### 4. /bluelm/chat (文本生成/聊天)

**请求方法**: POST  
**Content-Type**: application/json

**请求参数**:
```json
{
  "message": "你好，请介绍一下自己",           // 用户消息(必需)
  "session_id": "session_123",              // 会话ID(可选)
  "history_messages": [                      // 历史消息(可选)
    {
      "role": "user",
      "content": "之前的消息"
    }
  ]
}
```

**响应**:
```json
{
  "success": true,
  "message": "Chat completed successfully",
  "timestamp": "2024-01-01 12:00:00",
  "session_id": "session_123",
  "data": {
    "reply": "你好！我是蓝心大模型...",
    "role": "assistant",
    "messages": [/* 完整消息历史 */]
  }
}
```

**使用示例**:
```bash
# curl示例 - 新对话
curl -X POST http://127.0.0.1:8888/bluelm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下自己"
  }'

# curl示例 - 继续对话
curl -X POST http://127.0.0.1:8888/bluelm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我刚刚问了什么？",
    "session_id": "session_123",
    "history_messages": [
      {"role": "user", "content": "你好，请介绍一下自己"},
      {"role": "assistant", "content": "你好！我是蓝心大模型..."}
    ]
  }'
```

```javascript
// JavaScript示例 - 新对话
fetch('http://127.0.0.1:8888/bluelm/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    message: '你好，请介绍一下自己'
  })
})
.then(response => response.json())
.then(data => {
  console.log('AI回复:', data.data.reply);
  console.log('会话ID:', data.session_id);
});

// JavaScript示例 - 继续对话
fetch('http://127.0.0.1:8888/bluelm/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    message: '我刚刚问了什么？',
    session_id: 'session_123',
    history_messages: [
      {role: 'user', content: '你好，请介绍一下自己'},
      {role: 'assistant', content: '你好！我是蓝心大模型...'}
    ]
  })
})
.then(response => response.json())
.then(data => {
  console.log('AI回复:', data.data.reply);
});
```

## 环境要求

### 长语音转写whisperX

- 需要Python环境

​		版本python = ">=3.9, <3.13"

- 需要一个HuggingFace的token，读取权限即可

​		同时要先获取[segmentation-3.0](https://huggingface.co/pyannote/segmentation-3.0)和[speaker-diarization-3.1](https://huggingface.co/pyannote/speaker-diarization-3.1)的授权

### 蓝心大模型服务

- 需要配置环境变量:
  - `APPID`: 蓝心大模型应用ID
  - `APPKEY`: 蓝心大模型应用密钥

## 错误处理

所有接口在出错时都会返回相应的HTTP状态码和错误信息:

- `400 Bad Request`: 请求参数错误
- `500 Internal Server Error`: 服务器内部错误

错误响应格式:
```json
{
  "error": "错误类型",
  "message": "详细错误信息"
}
```

