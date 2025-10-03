package domain

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a background task
type Task struct {
	ID        string
	Type      string
	Data      map[string]interface{}
	CreatedAt time.Time
}

// TaskManager manages background tasks
type TaskManager struct {
	tasks     chan Task
	workers   int
	handlers  map[string]TaskHandler
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// TaskHandler handles a specific type of task
type TaskHandler interface {
	Handle(ctx context.Context, task Task) error
}

// NewTaskManager creates a new task manager
func NewTaskManager(workers int) *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskManager{
		tasks:    make(chan Task, 100), // Buffer for 100 tasks
		workers:  workers,
		handlers: make(map[string]TaskHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RegisterHandler registers a task handler
func (tm *TaskManager) RegisterHandler(taskType string, handler TaskHandler) {
	tm.handlers[taskType] = handler
}

// Start starts the task manager
func (tm *TaskManager) Start() {
	for i := 0; i < tm.workers; i++ {
		tm.wg.Add(1)
		go tm.worker(i)
	}
	log.Printf("Task manager started with %d workers", tm.workers)
}

// Stop stops the task manager
func (tm *TaskManager) Stop() {
	tm.cancel()
	close(tm.tasks)
	tm.wg.Wait()
	log.Println("Task manager stopped")
}

// Enqueue adds a task to the queue
func (tm *TaskManager) Enqueue(task Task) error {
	select {
	case tm.tasks <- task:
		log.Printf("Task %s of type %s enqueued", task.ID, task.Type)
		return nil
	case <-tm.ctx.Done():
		return tm.ctx.Err()
	default:
		log.Printf("Task queue full, dropping task %s", task.ID)
		return nil // Don't block on full queue
	}
}

// worker processes tasks from the queue
func (tm *TaskManager) worker(workerID int) {
	defer tm.wg.Done()
	
	log.Printf("Worker %d started", workerID)
	
	for {
		select {
		case task, ok := <-tm.tasks:
			if !ok {
				log.Printf("Worker %d stopping (channel closed)", workerID)
				return
			}
			
			tm.processTask(workerID, task)
			
		case <-tm.ctx.Done():
			log.Printf("Worker %d stopping (context cancelled)", workerID)
			return
		}
	}
}

// processTask processes a single task
func (tm *TaskManager) processTask(workerID int, task Task) {
	log.Printf("Worker %d processing task %s of type %s", workerID, task.ID, task.Type)
	
	handler, exists := tm.handlers[task.Type]
	if !exists {
		log.Printf("Worker %d: No handler for task type %s", workerID, task.Type)
		return
	}
	
	start := time.Now()
	err := handler.Handle(tm.ctx, task)
	duration := time.Since(start)
	
	if err != nil {
		log.Printf("Worker %d: Task %s failed after %v: %v", workerID, task.ID, duration, err)
	} else {
		log.Printf("Worker %d: Task %s completed in %v", workerID, task.ID, duration)
	}
}