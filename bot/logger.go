package bot

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Unanoc/InstaFollower/pkg"
)

// Logger ensures th logging
type Logger struct {
	LogPath string
}

// CreateLogger return th instance of Logger
func CreateLogger(path string) *Logger {
	return &Logger{LogPath: path}
}

// Log writes logs to file. If the file does not exist, creates it
func (l *Logger) Log(msgType string, args ...string) {
	logFileName := time.Now().Format("2006-01-02")
	logfilePath := pkg.GetCorrectPath(l.LogPath, logFileName, ".log")

	f, err := os.OpenFile(logfilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	logContent := strings.Join(args, " | ")
	log.Println(logContent)
	logger := log.New(f, msgType, log.LstdFlags)
	logger.Print(logContent)
}
