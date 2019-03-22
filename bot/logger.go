package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Logger ...
type Logger struct {
	LogPath string
}

// CreateLogger ...
func CreateLogger(path string) *Logger {
	return &Logger{LogPath: path}
}

// Log ...
func (l *Logger) Log(username, text string) {
	logfilePath := getLogFilePath(l.LogPath)

	f, err := os.OpenFile(logfilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	logger := log.New(f, "[MESSAGE] ", log.LstdFlags)
	logger.Printf("[%s] %s", username, text)
}

// LogInfo ...
func (l *Logger) LogInfo(text string) {
	logfilePath := getLogFilePath(l.LogPath)

	f, err := os.OpenFile(logfilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	logger := log.New(f, "[INFO] ", log.LstdFlags)
	logger.Print(text)
}

func getLogFilePath(logdir string) string {
	if len(logdir) == 0 {
		return ""
	}

	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		os.Mkdir(logdir, 0700)
	}

	date := time.Now().Format("2006-01-02")
	if strings.LastIndex(logdir, "/") == len(logdir)-1 {
		return fmt.Sprintf("%s%s.log", logdir, date)
	}
	return fmt.Sprintf("%s/%s.log", logdir, date)
}
