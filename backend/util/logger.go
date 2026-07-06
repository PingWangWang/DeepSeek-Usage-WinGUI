package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 全局日志实例
var Logger *LoggerImpl

// LoggerImpl 日志实现
type LoggerImpl struct {
	file *os.File
}

func init() {
	Logger = &LoggerImpl{}
	Logger.init()
}

// init 初始化日志
func (l *LoggerImpl) init() {
	logDir := filepath.Join(os.Getenv("APPDATA"), "DeepSeek-Usage", "logs")
	os.MkdirAll(logDir, 0755)

	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}
	l.file = f
}

// Info 输出信息级别日志
func (l *LoggerImpl) Info(msg string, keysAndValues ...interface{}) {
	l.log("INFO", msg, keysAndValues...)
}

// Warn 输出警告级别日志
func (l *LoggerImpl) Warn(msg string, keysAndValues ...interface{}) {
	l.log("WARN", msg, keysAndValues...)
}

// Error 输出错误级别日志
func (l *LoggerImpl) Error(msg string, keysAndValues ...interface{}) {
	l.log("ERROR", msg, keysAndValues...)
}

// Debug 输出调试级别日志
func (l *LoggerImpl) Debug(msg string, keysAndValues ...interface{}) {
	l.log("DEBUG", msg, keysAndValues...)
}

// log 内部日志输出方法
func (l *LoggerImpl) log(level string, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s %s", timestamp, level, msg)

	// 添加键值对
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			logMsg += fmt.Sprintf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
		}
	}

	// 输出到控制台
	fmt.Println(logMsg)

	// 输出到文件
	if l.file != nil {
		fmt.Fprintln(l.file, logMsg)
	}
}

// Close 关闭日志文件
func (l *LoggerImpl) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
