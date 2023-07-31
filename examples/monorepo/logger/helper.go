package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var (
	instance *log.Logger
	file     *os.File
)

// Default returns the default logger.
func Default() *log.Logger {
	return instance
}

// Initialize creates the default logger and
// opens the file to keep log if enabled.
func Initialize() error {
	var writer io.Writer = os.Stdout
	if defaultCfg.file.enabled {
		filename := path.Join(defaultCfg.file.dir, defaultCfg.file.name)
		file, err := os.OpenFile(filename,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}

		writer = io.MultiWriter(os.Stdout, file)
	}

	instance = log.New(writer, "", log.Ldate|log.Ltime)

	return nil
}

// Finalize closes the log file.
func Finalize() {
	if file != nil {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
