package world

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/msmhq/msm/pkg/screen"
)

func (w *World) ToRAM(user string) error {
	if !w.GlobalCfg.RamdiskStorageEnabled {
		return fmt.Errorf("ramdisk storage is not enabled")
	}

	if err := os.MkdirAll(w.RAMPath, 0755); err != nil {
		return fmt.Errorf("failed to create RAM directory: %w", err)
	}

	flagPath := w.FlagPath()
	rsyncCmd := fmt.Sprintf("rsync -rt --exclude '%s' '%s/' '%s'",
		filepath.Base(flagPath), w.Path, w.RAMPath)

	return screen.RunAsUser(user, rsyncCmd)
}

func (w *World) ToDisk(user string) error {
	if !w.InRAM {
		return nil
	}

	flagPath := w.FlagPath()
	rsyncCmd := fmt.Sprintf("rsync -rt --exclude '%s' '%s/' '%s'",
		filepath.Base(flagPath), w.RAMPath, w.Path)

	return screen.RunAsUser(user, rsyncCmd)
}

func (w *World) EnableRAM(user string) error {
	if !w.GlobalCfg.RamdiskStorageEnabled {
		return fmt.Errorf("ramdisk storage is not enabled in config")
	}

	if w.InRAM {
		return fmt.Errorf("world %q is already in RAM", w.Name)
	}

	flagPath := w.FlagPath()

	fmt.Printf("Adding RAM flag to world %q... ", w.Name)
	if err := touchFile(flagPath); err != nil {
		return err
	}
	fmt.Println("Done.")

	fmt.Printf("Copying world %q to RAM... ", w.Name)
	w.RAMPath = filepath.Join(w.GlobalCfg.RamdiskStoragePath, w.ServerName, w.Name)
	if err := w.ToRAM(user); err != nil {
		return err
	}
	fmt.Println("Done.")

	fmt.Printf("Updating allowed_symlinks.txt... ")
	if err := EnsureAllowedSymlinks(w.ServerPath, w.GlobalCfg.RamdiskStoragePath); err != nil {
		return err
	}
	fmt.Println("Done.")

	w.InRAM = true
	fmt.Println("RAM enabled. Changes will take effect after server is restarted.")
	return nil
}

func (w *World) DisableRAM(user string) error {
	if !w.InRAM {
		return fmt.Errorf("world %q is not in RAM", w.Name)
	}

	flagPath := w.FlagPath()

	fmt.Printf("Synchronising world %q to disk... ", w.Name)
	if err := w.ToDisk(user); err != nil {
		return err
	}
	fmt.Println("Done.")

	fmt.Printf("Removing RAM flag from world %q... ", w.Name)
	if err := os.Remove(flagPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Println("Done.")

	fmt.Printf("Removing world %q from RAM... ", w.Name)
	if err := os.RemoveAll(w.RAMPath); err != nil {
		return err
	}
	fmt.Println("Done.")

	fmt.Printf("Updating allowed_symlinks.txt... ")
	if err := RemoveAllowedSymlink(w.ServerPath, w.GlobalCfg.RamdiskStoragePath); err != nil {
		return err
	}
	fmt.Println("Done.")

	w.InRAM = false
	fmt.Println("RAM disabled. Changes will take effect after server is restarted.")
	return nil
}

func (w *World) SetupRAMSymlink(worldStoragePath string) error {
	if !w.InRAM || !w.GlobalCfg.RamdiskStorageEnabled {
		return nil
	}

	if err := EnsureAllowedSymlinks(w.ServerPath, w.GlobalCfg.RamdiskStoragePath); err != nil {
		return fmt.Errorf("failed to update allowed_symlinks.txt: %w", err)
	}

	symlinkPath := filepath.Join(w.ServerPath, w.Name)

	linkTarget, err := os.Readlink(symlinkPath)
	if err == nil && linkTarget == w.RAMPath {
		return nil
	}

	info, err := os.Lstat(symlinkPath)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(symlinkPath); err != nil {
				return fmt.Errorf("failed to remove existing symlink: %w", err)
			}
		} else if info.IsDir() {
			return fmt.Errorf("world directory exists at %s - expected symlink to worldstorage", symlinkPath)
		}
	}

	if err := os.Symlink(w.RAMPath, symlinkPath); err != nil {
		return fmt.Errorf("failed to create RAM symlink: %w", err)
	}

	return nil
}

func touchFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return file.Close()
}

func EnsureAllowedSymlinks(serverPath, ramdiskPath string) error {
	allowedFile := filepath.Join(serverPath, "allowed_symlinks.txt")
	entry := "prefix" + ramdiskPath

	existingLines := []string{}
	if file, err := os.Open(allowedFile); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			existingLines = append(existingLines, line)
			if line == entry {
				file.Close()
				return nil
			}
		}
		file.Close()
	}

	existingLines = append(existingLines, entry)

	file, err := os.Create(allowedFile)
	if err != nil {
		return fmt.Errorf("failed to create allowed_symlinks.txt: %w", err)
	}
	defer file.Close()

	for _, line := range existingLines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("failed to write allowed_symlinks.txt: %w", err)
		}
	}

	return nil
}

func RemoveAllowedSymlink(serverPath, ramdiskPath string) error {
	allowedFile := filepath.Join(serverPath, "allowed_symlinks.txt")
	entry := "prefix" + ramdiskPath

	file, err := os.Open(allowedFile)
	if err != nil {
		return nil
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, entry) {
			lines = append(lines, line)
		}
	}
	file.Close()

	if len(lines) == 0 {
		return os.Remove(allowedFile)
	}

	outFile, err := os.Create(allowedFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, line := range lines {
		fmt.Fprintln(outFile, line)
	}

	return nil
}
