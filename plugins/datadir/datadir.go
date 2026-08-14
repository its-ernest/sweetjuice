package datadir

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type DataDirPlugin struct {
	app *core.Application
}

type AppDirs struct {
	Files         string `json:"files"`
	Cache         string `json:"cache"`
	ExternalFiles string `json:"external_files"`
	ExternalCache string `json:"external_cache"`
}

// FileResult represents the result of a file operation.
type FileResult struct {
	Success bool   `json:"success"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewPlugin creates a new instance of the DataDir plugin.
func NewPlugin() *DataDirPlugin {
	return &DataDirPlugin{}
}

// Name returns the plugin name.
func (p *DataDirPlugin) Name() string {
	return "datadir"
}

// Init initializes the plugin.
func (p *DataDirPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// GetDirs returns the standard application directories.
func (p *DataDirPlugin) GetDirs() (AppDirs, error) {
	var dirs AppDirs
	result := core.CallNativePlatform("datadir:getDirs", "{}")

	if err := parseResultError(result); err != nil {
		return dirs, err
	}

	if err := json.Unmarshal([]byte(result), &dirs); err != nil {
		return dirs, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return dirs, nil
}

// ReadFile reads a file from the app's internal files directory.
func (p *DataDirPlugin) ReadFile(path string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"path": path,
	})
	result := core.CallNativePlatform("datadir:readFile", string(payload))

	if err := parseResultError(result); err != nil {
		return "", err
	}

	var fileResult FileResult
	if err := json.Unmarshal([]byte(result), &fileResult); err != nil {
		return "", fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	if !fileResult.Success {
		return "", fmt.Errorf("readFile failed: %s", fileResult.Error)
	}

	return fileResult.Content, nil
}

// WriteFile writes content to a file in the app's internal files directory.
func (p *DataDirPlugin) WriteFile(path string, content string) error {
	payload, _ := json.Marshal(map[string]string{
		"path":    path,
		"content": content,
	})
	result := core.CallNativePlatform("datadir:writeFile", string(payload))

	if err := parseResultError(result); err != nil {
		return err
	}

	var fileResult FileResult
	if err := json.Unmarshal([]byte(result), &fileResult); err != nil {
		return fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	if !fileResult.Success {
		return fmt.Errorf("writeFile failed: %s", fileResult.Error)
	}

	return nil
}

// FileExists checks whether a file exists in the app's internal files directory.
func (p *DataDirPlugin) FileExists(path string) (bool, error) {
	payload, _ := json.Marshal(map[string]string{
		"path": path,
	})
	result := core.CallNativePlatform("datadir:fileExists", string(payload))

	if err := parseResultError(result); err != nil {
		return false, err
	}

	var fileResult FileResult
	if err := json.Unmarshal([]byte(result), &fileResult); err != nil {
		return false, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return fileResult.Success, nil
}

// DeleteFile deletes a file from the app's internal files directory.
func (p *DataDirPlugin) DeleteFile(path string) error {
	payload, _ := json.Marshal(map[string]string{
		"path": path,
	})
	result := core.CallNativePlatform("datadir:deleteFile", string(payload))

	if err := parseResultError(result); err != nil {
		return err
	}

	var fileResult FileResult
	if err := json.Unmarshal([]byte(result), &fileResult); err != nil {
		return fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	if !fileResult.Success {
		return fmt.Errorf("deleteFile failed: %s", fileResult.Error)
	}

	return nil
}

func parseResultError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"]; ok {
		return fmt.Errorf("native error: %v", errMsg)
	}
	return nil
}
