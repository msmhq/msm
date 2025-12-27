package world

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/msmhq/msm/pkg/screen"
)

func (w *World) Backup(backupPath, user string) error {
	if w.GlobalCfg.WorldArchiveEnabled {
		if err := w.backupZip(backupPath, user); err != nil {
			return fmt.Errorf("zip backup failed: %w", err)
		}
	}

	if w.GlobalCfg.RdiffBackupEnabled {
		if err := w.backupRdiff(user); err != nil {
			return fmt.Errorf("rdiff backup failed: %w", err)
		}
	}

	if w.GlobalCfg.RsyncBackupEnabled {
		if err := w.backupRsync(user); err != nil {
			return fmt.Errorf("rsync backup failed: %w", err)
		}
	}

	return nil
}

func (w *World) backupZip(backupPath, user string) error {
	if backupPath == "" {
		backupPath = filepath.Join(w.GlobalCfg.WorldArchivePath, w.ServerName, w.Name)
	}

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	zipPath := filepath.Join(backupPath, timestamp+".zip")

	fmt.Printf("Backing up world %q... ", w.Name)

	sourcePath := w.Path
	if w.InRAM && w.RAMPath != "" {
		sourcePath = w.RAMPath
	}

	if err := createZip(zipPath, sourcePath); err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}

func (w *World) backupRdiff(user string) error {
	backupPath := filepath.Join(w.GlobalCfg.WorldRdiffPath, w.ServerName, w.Name)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create rdiff backup directory: %w", err)
	}

	fmt.Printf("rdiff-backup world %q... ", w.Name)

	sourcePath := w.Path
	if w.InRAM && w.RAMPath != "" {
		sourcePath = w.RAMPath
	}

	rdiffCmd := fmt.Sprintf("nice -n %d rdiff-backup '%s' '%s' && nice -n %d rdiff-backup --remove-older-than %dD --force '%s'",
		w.GlobalCfg.RdiffBackupNice, sourcePath, backupPath,
		w.GlobalCfg.RdiffBackupNice, w.GlobalCfg.RdiffBackupRotation, backupPath)

	if err := screen.RunAsUser(user, rdiffCmd); err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}

func (w *World) backupRsync(user string) error {
	backupPath := filepath.Join(w.GlobalCfg.WorldRsyncPath, w.ServerName, w.Name)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create rsync backup directory: %w", err)
	}

	fmt.Printf("rsync-backup world %q... ", w.Name)

	sourcePath := w.Path
	if w.InRAM && w.RAMPath != "" {
		sourcePath = w.RAMPath
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	latestPath := filepath.Join(backupPath, "latest")
	destPath := filepath.Join(backupPath, timestamp)

	rsyncCmd := fmt.Sprintf("rsync -aH --link-dest='%s' '%s' '%s' && rm -f '%s' && ln -s '%s' '%s'",
		latestPath, sourcePath, destPath, latestPath, timestamp, latestPath)

	if err := screen.RunAsUser(user, rsyncCmd); err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}

func createZip(zipPath, sourceDir string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	baseDir := filepath.Base(sourceDir)

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(baseDir, relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func BackupServer(serverPath, serverName, backupArchivePath, user string, followSymlinks bool) error {
	if err := os.MkdirAll(backupArchivePath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	zipPath := filepath.Join(backupArchivePath, serverName, timestamp+".zip")

	if err := os.MkdirAll(filepath.Dir(zipPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	fmt.Printf("Backing up server %q... ", serverName)

	if err := createZip(zipPath, serverPath); err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}
