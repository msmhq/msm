package screen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Session struct {
	Name string
}

func New(name string) *Session {
	return &Session{Name: name}
}

func (s *Session) IsRunning() bool {
	return s.IsRunningAsUser("")
}

func (s *Session) IsRunningAsUser(user string) bool {
	var cmd *exec.Cmd
	if user != "" && user != currentUser() && IsRoot() {
		screenCmd := fmt.Sprintf("screen -ls %s", s.Name)
		cmd = exec.Command("su", "-", user, "-s", "/bin/bash", "-c", screenCmd)
	} else {
		cmd = exec.Command("screen", "-ls", s.Name)
	}
	output, _ := cmd.Output()
	return strings.Contains(string(output), s.Name)
}

func (s *Session) Start(workDir, command string, user string) error {
	if s.IsRunningAsUser(user) {
		return fmt.Errorf("session %q is already running", s.Name)
	}

	args := []string{"-dmS", s.Name, "bash", "-c", command}

	cmd := exec.Command("screen", args...)
	cmd.Dir = workDir

	if user != "" && user != currentUser() {
		return s.startAsUser(workDir, command, user)
	}

	return cmd.Run()
}

func (s *Session) startAsUser(workDir, command, user string) error {
	screenCmd := fmt.Sprintf("cd %q && screen -dmS %s bash -c %q", workDir, s.Name, command)
	cmd := exec.Command("su", "-", user, "-s", "/bin/bash", "-c", screenCmd)
	return cmd.Run()
}

func (s *Session) SendCommand(command string) error {
	if !s.IsRunning() {
		return fmt.Errorf("session %q is not running", s.Name)
	}

	cmd := exec.Command("screen", "-S", s.Name, "-p", "0", "-X", "stuff", command+"\n")
	return cmd.Run()
}

func (s *Session) SendCommandAsUser(command, user string) error {
	if user != "" && user != currentUser() {
		screenCmd := fmt.Sprintf("screen -S %s -p 0 -X stuff '%s\n'", s.Name, command)
		cmd := exec.Command("su", "-", user, "-s", "/bin/bash", "-c", screenCmd)
		return cmd.Run()
	}
	return s.SendCommand(command)
}

func (s *Session) Attach() error {
	if !s.IsRunning() {
		return fmt.Errorf("session %q is not running", s.Name)
	}

	cmd := exec.Command("screen", "-r", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Session) AttachAsUser(user string) error {
	if user != "" && user != currentUser() {
		screenCmd := fmt.Sprintf("screen -r %s", s.Name)
		cmd := exec.Command("su", "-", user, "-s", "/bin/bash", "-c", screenCmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return s.Attach()
}

func (s *Session) Kill() error {
	if !s.IsRunning() {
		return nil
	}

	cmd := exec.Command("screen", "-S", s.Name, "-X", "quit")
	return cmd.Run()
}

func (s *Session) GetPID() (int, error) {
	cmd := exec.Command("screen", "-ls")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	for _, line := range strings.Split(out.String(), "\n") {
		if strings.Contains(line, s.Name) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				pidPart := strings.Split(fields[0], ".")[0]
				pid, err := strconv.Atoi(pidPart)
				if err != nil {
					continue
				}
				return pid, nil
			}
		}
	}

	return 0, fmt.Errorf("session %q not found", s.Name)
}

func currentUser() string {
	return os.Getenv("USER")
}

func RunAsUser(user, command string) error {
	if user == "" || user == currentUser() {
		cmd := exec.Command("bash", "-c", command)
		return cmd.Run()
	}

	cmd := exec.Command("su", "-", user, "-s", "/bin/bash", "-c", command)
	return cmd.Run()
}

func RunAsUserWithOutput(user, command string) (string, error) {
	var cmd *exec.Cmd

	if user == "" || user == currentUser() {
		cmd = exec.Command("bash", "-c", command)
	} else {
		cmd = exec.Command("su", "-", user, "-s", "/bin/bash", "-c", command)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}

func IsRoot() bool {
	return syscall.Getuid() == 0
}

func CurrentUser() string {
	return currentUser()
}

func CanManageUser(targetUser string) bool {
	if IsRoot() {
		return true
	}
	if targetUser == "" {
		return true
	}
	return currentUser() == targetUser
}
