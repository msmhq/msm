package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func Roll(logPath, archivePath, serverName string) error {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil
	}

	destDir := filepath.Join(archivePath, serverName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	destPath := filepath.Join(destDir, timestamp+".log.gz")

	fmt.Printf("Rolling log to %s... ", destPath)

	srcFile, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to compress log: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to finalize compression: %w", err)
	}

	srcFile.Close()

	if err := os.Truncate(logPath, 0); err != nil {
		return fmt.Errorf("failed to truncate log file: %w", err)
	}

	fmt.Println("Done.")
	return nil
}
