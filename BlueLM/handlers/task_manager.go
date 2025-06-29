package handlers

import (
	"sync"
	"time"
)

// TaskStatus 定义任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

// TaskInfo 存储任务信息
type TaskInfo struct {
	TaskID    string     `json:"task_id"`
	Status    TaskStatus `json:"status"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	Filename  string     `json:"filename"`
	FilePath  string     `json:"-"` // 不在 JSON 中暴露文件路径
}

// TaskManager 管理转录任务
type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*TaskInfo
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*TaskInfo),
	}
}

// CreateTask 创建新任务
func (tm *TaskManager) CreateTask(taskID, filename string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	tm.tasks[taskID] = &TaskInfo{
		TaskID:    taskID,
		Status:    TaskStatusPending,
		Message:   "Task created, waiting to start",
		CreatedAt: time.Now(),
		Filename:  filename,
	}
}

// UpdateTaskStatus 更新任务状态
func (tm *TaskManager) UpdateTaskStatus(taskID string, status TaskStatus, message string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if task, exists := tm.tasks[taskID]; exists {
		task.Status = status
		task.Message = message
	}
}

// SetTaskFilePath 设置任务的文件路径
func (tm *TaskManager) SetTaskFilePath(taskID, filePath string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if task, exists := tm.tasks[taskID]; exists {
		task.FilePath = filePath
	}
}

// GetTask 获取任务信息
func (tm *TaskManager) GetTask(taskID string) (*TaskInfo, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, false
	}
	
	// 返回副本以避免并发问题
	taskCopy := *task
	return &taskCopy, true
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() []*TaskInfo {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	tasks := make([]*TaskInfo, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		// 返回副本以避免并发问题
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}
	
	return tasks
}

// 全局任务管理器实例
var GlobalTaskManager = NewTaskManager()