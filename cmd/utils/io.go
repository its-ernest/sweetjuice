// Package utils provides shared cross-platform utility functions for file I/O,
// networking, and process execution.
package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists returns true if the specified path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// DirEmpty returns true if the directory is empty or cannot be opened.
func DirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	return err == io.EOF
}

// IsWindowsHost returns true if the current operating system is Windows.
func IsWindowsHost() bool {
	return runtime.GOOS == "windows"
}

// CommandExists returns true if the specified command is available in the system PATH.
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// EnsureGoMobileTools checks if gomobile and gobind are installed, and installs them if missing.
func EnsureGoMobileTools() {
	if !CommandExists("gomobile") || !CommandExists("gobind") {
		Info("Go Mobile build tools missing. Installing...")
		RunCmd("go", "install", "golang.org/x/mobile/cmd/gomobile@latest")
		RunCmd("go", "install", "golang.org/x/mobile/cmd/gobind@latest")

		Info("Initializing Go Mobile environment...")
		RunCmd("gomobile", "init")
	}
}

// RunCmd executes an external command and pipes its output directly to the current process's stdout and stderr.
func RunCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Fatal(fmt.Sprintf("Command failed: %s", name), err)
	}
}

// ReplaceInFiles recursively searches for a string in all files within a directory and replaces it with another.
func ReplaceInFiles(root, old, new string, extension string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, extension) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		if strings.Contains(content, old) {
			newContent := strings.ReplaceAll(content, old, new)
			err = os.WriteFile(path, []byte(newContent), info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// BuildFrontend executes the frontend build command defined in the config.ini file.
func BuildFrontend() {
	config := LoadConfig()
	buildCommand := config.GetOrDefault("build", "frontend_build_command", "")
	if buildCommand == "" {
		return
	}

	frontendDir := config.GetOrDefault("build", "frontend_dir", "frontend")
	if !DirExists(frontendDir) {
		return
	}

	Info("Building frontend with command: " + buildCommand)

	origWd, _ := os.Getwd()
	_ = os.Chdir(frontendDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Execute command via shell to support && and complex syntax
	var cmd *exec.Cmd
	if IsWindowsHost() {
		cmd = exec.Command("cmd", "/C", buildCommand)
	} else {
		cmd = exec.Command("sh", "-c", buildCommand)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Warn("Frontend build failed: " + err.Error())
	}
}

// DownloadFile retrieves a file from the given URL and saves it to the specified destination path.
func DownloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// UnzipTarget extracts a ZIP archive from src to the dest directory.
func UnzipTarget(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyFile copies a single file from src to dest.
func CopyFile(src, dest string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	info, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer d.Close()

	if _, err := io.Copy(d, s); err != nil {
		return err
	}
	return nil
}

// MoveFile renames a file, or copies and deletes it if a rename fails (e.g., cross-device).
func MoveFile(src, dest string) error {
	err := os.Rename(src, dest)
	if err == nil {
		return nil
	}

	// Fallback to copy and delete
	if err := CopyFile(src, dest); err != nil {
		return err
	}
	return os.Remove(src)
}

// CopyDirectory recursively copies a directory and its contents from scrDir to destDir, excluding .git.
func CopyDirectory(scrDir, destDir string) error {
	return filepath.Walk(scrDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(scrDir, path)
		if relPath == ".git" || strings.HasPrefix(relPath, ".git/") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		targetPath := filepath.Join(destDir, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// CopyAppAssets copies everything from the app_assets directory into the given target directory.
// If app_assets does not exist, it returns nil without error.
func CopyAppAssets(targetDir string) error {
	srcDir := "app_assets"
	if !DirExists(srcDir) {
		return nil
	}

	if !DirExists(targetDir) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target assets directory %s: %w", targetDir, err)
		}
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read app_assets directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			if err := CopyDirectory(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy app_assets directory %s: %w", entry.Name(), err)
			}
		} else {
			if err := CopyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy app_assets file %s: %w", entry.Name(), err)
			}
		}
	}

	Info(fmt.Sprintf("Copied app_assets into %s", targetDir))
	return nil
}
