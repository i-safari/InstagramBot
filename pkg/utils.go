package pkg

import (
	"fmt"
	"os"
	"strings"
)

// GetCorrectPath ...
func GetCorrectPath(pathToDir, fileName, ext string) string {
	if len(pathToDir) == 0 {
		return ""
	}

	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		os.Mkdir(pathToDir, 0700)
	}

	if strings.LastIndex(pathToDir, "/") == len(pathToDir)-1 {
		return fmt.Sprintf("%s%s%s", pathToDir, fileName, ext)
	}
	return fmt.Sprintf("%s/%s%s", pathToDir, fileName, ext)
}
