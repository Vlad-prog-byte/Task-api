package logger

import (
	"fmt"
	"log"
	"os"
	"task-api/internals/models"
	"time"
	"sync"
)

type TaskLogType string

const (
	LogCreate TaskLogType = "CREATE"
	LogDelete TaskLogType = "DELETE"
	LogUpdate TaskLogType = "UPDATE"
)

type LogEvent struct {
	Type TaskLogType
	Task models.Task
}

type TaskLogger interface {
	Log(event LogEvent)
	Close()
}

type taskLogger struct {
	logChan chan LogEvent
	file *os.File
	mu sync.Mutex
	wg sync.WaitGroup
}

func (logger *taskLogger) Log(event LogEvent) {
	logger.logChan <- event
}

func (logger *taskLogger) Close() {
	close(logger.logChan)
	logger.wg.Wait()
	logger.file.Close()
}

func NewTaskLogger() TaskLogger {
	f, err := os.OpenFile("logger.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := taskLogger{
		logChan: make(chan LogEvent),
		file: f,
	}
	logger.wg.Add(1)

	go func() {
		defer logger.wg.Done()
		for event := range logger.logChan {
			currentTime := time.Now()
			logMessage := fmt.Sprintf("%s [%s] Task ID: %d | Title: %s\n", currentTime.Format("2006-01-02 15:04:05"), event.Type, event.Task.ID, event.Task.Title)
			fmt.Println(logMessage)
			logger.mu.Lock()
			f.Write([]byte(logMessage))
			logger.mu.Unlock()
		}
	}()
	return &logger
}


