package app

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogToFile adds writer to the logger output destination, that writes all logs to the specified file. Truncates file at start
func LogToFile(file string) (io.Writer, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open logs file: %v", err)
	}
	
	customWriter := io.MultiWriter(f, os.Stdout)
	log.SetOutput(customWriter)
	log.SetFlags(log.LstdFlags|log.Ltime|log.Lshortfile)
	log.SetPrefix("[PAPRIKA]: ")
	return customWriter, nil
}