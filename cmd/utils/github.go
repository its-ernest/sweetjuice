package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func EnsureGHCLILoggedIn() {
	if !CommandExists("gh") {
		Fatal("GitHub CLI (gh) missing", fmt.Errorf("please install the GitHub CLI: https://cli.github.com"))
	}
}

func ForkAndCloneRepo(repoOwner, repoName, targetLocalPath string) {
	if DirExists(targetLocalPath) {
		return
	}

	Info("Forking and cloning " + repoOwner + "/" + repoName + " to " + targetLocalPath + "...")

	_ = exec.Command("gh", "repo", "fork", repoOwner+"/"+repoName, "--clone=false").Run()

	out, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		Fatal("Failed to get GitHub user info", err)
	}
	user := strings.TrimSpace(string(out))

	RunCmd("gh", "repo", "clone", user+"/"+repoName, targetLocalPath)
}

func WaitForActionFinish(repoPath string) string {
	Info("Waiting for GitHub Action to start and complete...")

	origWd, _ := os.Getwd()
	_ = os.Chdir(repoPath)
	defer func() { _ = os.Chdir(origWd) }()

	for i := 0; i < 60; i++ {
		time.Sleep(10 * time.Second)

		out, err := exec.Command("gh", "run", "list", "--limit", "1", "--json", "status,conclusion,databaseId").Output()
		if err != nil {
			continue
		}

		output := string(out)
		if strings.Contains(output, "\"status\":\"completed\"") {
			if strings.Contains(output, "\"conclusion\":\"success\"") {
				parts := strings.Split(output, "\"databaseId\":")
				if len(parts) > 1 {
					idPart := strings.Split(parts[1], "}")[0]
					idPart = strings.Split(idPart, ",")[0]
					idPart = strings.Split(idPart, "]")[0]
					return strings.Trim(strings.TrimSpace(idPart), "\"")
				}
			} else if strings.Contains(output, "\"conclusion\":\"failure\"") {
				Fatal("GitHub Action failed", fmt.Errorf("check 'gh run view' for details"))
			}
		}
		Print(".")
	}

	Fatal("Timeout waiting for GitHub Action", fmt.Errorf("the build is taking too long"))
	return ""
}

func DownloadArtifact(repoPath, runID, artifactName, destPath string) {
	Println()
	Info("Downloading artifact " + artifactName + "...")

	origWd, _ := os.Getwd()
	_ = os.Chdir(repoPath)
	defer func() { _ = os.Chdir(origWd) }()

	RunCmd("gh", "run", "download", runID, "--name", artifactName, "--dir", destPath)
}
