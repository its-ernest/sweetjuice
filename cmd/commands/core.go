package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/android"
	"github.com/sweet-juice/sweetjuice/cmd/ios"
	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

const (
	Version = "v1.0.0"
)

func ShowUsage() {
	utils.Println("Sweet Juice Toolchain CLI (juice)")
	utils.Println("Usage:")
	utils.Println("  juice --new <project_name>        Create a fresh project from the template")
	utils.Println("  juice --refresh <platform>        Run platform sync: 'android' or 'ios'")
	utils.Println("  juice --build <platform> <mode>   Compile binaries: 'debug' (APK/IPA), 'release' (APK/IPA), or 'bundle' (AAB)")
	utils.Println("  juice --run <platform>            Compile, install, and execute application via ADB or xtool")
	utils.Println("  juice --run --force <platform>    Run application without rebuilding bindings (assumes xcframework/aar present)")
	utils.Println("  juice --run-cross <platform>      Cloud build (GitHub Actions), install, and execute")
	utils.Println("  juice --setup cross               Setup GitHub Action based cross-compilation for iOS")
	utils.Println("  juice --add <plugin-url>          Install a native Go/Mobile plugin")
	utils.Println("  juice --remove <plugin-url>       Uninstall a native Go/Mobile plugin")
	os.Exit(1)
}

func CreateNewProject(targetDir string) {
	utils.Step("Creating Project: " + targetDir + " [" + Version + "]")
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		utils.Error("Directory '" + targetDir + "' already exists. Aborting.")
		os.Exit(1)
	}

	if _, err := exec.LookPath("go"); err != nil {
		utils.Error("Required system tool dependency 'go' is missing.")
		os.Exit(1)
	}

	utils.Info("Locating Sweet Juice core...")
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/sweet-juice/sweetjuice").Output()
	coreDir := strings.TrimSpace(string(out))

	if err != nil || coreDir == "" {
		utils.Info("Core not found in local context, checking module cache...")
		out, err = exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/sweet-juice/sweetjuice@latest").Output()
		coreDir = strings.TrimSpace(string(out))
	}

	if coreDir == "" {
		exePath, _ := os.Executable()
		repoCandidate := filepath.Join(filepath.Dir(exePath), "..", "..")
		if utils.FileExists(filepath.Join(repoCandidate, "go.mod")) {
			coreDir, _ = filepath.Abs(repoCandidate)
		}
	}

	if coreDir == "" {
		utils.Fatal("Failed to locate Sweet Juice core directory", fmt.Errorf("please ensure you have github.com/sweet-juice/sweetjuice installed or are running from the source repo"))
	}

	templatePath := filepath.Join(coreDir, "AppTemplate")

	if !utils.DirExists(templatePath) {
		utils.Fatal("Could not find AppTemplate in core module directory", fmt.Errorf("path missing: %s", templatePath))
	}

	utils.Info("Scaffolding project from local template...")
	if err := utils.CopyDirectory(templatePath, targetDir); err != nil {
		utils.Fatal("Failed to copy template to target directory", err)
	}

	_ = os.RemoveAll(filepath.Join(targetDir, "frontend", "node_modules"))

	targetGoMod := filepath.Join(targetDir, "go.mod")
	if utils.FileExists(targetGoMod) && coreDir != "" {
		utils.Info("Configuring project dependencies...")
		utils.RunCmd("go", "mod", "edit", "-dropreplace=github.com/sweet-juice/sweetjuice", targetGoMod)
		utils.RunCmd("go", "mod", "edit", "-replace=github.com/sweet-juice/sweetjuice="+coreDir, targetGoMod)
	}

	android.SetupAndroidLocalProperties(targetDir)

	origWd, _ := os.Getwd()
	_ = os.Chdir(targetDir)
	defer func() { _ = os.Chdir(origWd) }()

	if utils.DirExists("frontend") {
		utils.BuildFrontend()
	}

	utils.EnsureGoMobileTools()

	if utils.FileExists("go.mod") {
		utils.Info("Valid Go module context detected. Binding tool tracking dependencies...")
		utils.RunCmd("go", "mod", "tidy")
		if !utils.CommandExists("gobind") {
			utils.RunCmd("go", "get", "-tool", "golang.org/x/mobile/cmd/gobind")
		}
	}

	utils.Success("Setup complete! Your project is ready in ./" + targetDir)
}

func ExecuteRefresh(platform string) {
	if platform != "android" && platform != "ios" {
		utils.Error("Please specify a valid target platform: 'android' or 'ios'")
		os.Exit(1)
	}

	utils.Step("Executing Platform Refresh: " + platform)
	if platform == "android" {
		android.RefreshPipeline()
	} else {
		ios.RefreshPipeline()
	}
}

func ExecuteBuild(platform, mode string) {
	mode = strings.ToLower(mode)
	if mode != "debug" && mode != "release" && mode != "bundle" {
		utils.Error("Invalid build mode specified. Use 'debug' (APK/IPA), 'release' (APK/IPA), or 'bundle' (AAB).")
		os.Exit(1)
	}

	ExecuteRefresh(platform)

	if platform == "android" {
		android.BuildPipeline(mode)
	} else {
		ios.BuildPipeline(mode)
	}
}

func ExecuteRun(platform string, force bool) {
	if platform == "android" {
		ExecuteBuild("android", "debug")
		android.RunPipeline()
	} else if platform == "ios" {
		if force {
			ios.RunPipeline()
		} else {
			ExecuteBuild("ios", "debug")
			ios.RunPipeline()
		}
	} else {
		utils.Error("Invalid platform. Use 'android' or 'ios'.")
		os.Exit(1)
	}
}

func ExecuteRunCross(platform string) {
	if platform != "ios" {
		utils.Warn("Cross-build is currently only optimized for 'ios'. Running standard local build for others...")
		ExecuteRun(platform, false)
		return
	}

	config := utils.LoadConfig()
	githubUser := config.GetOrDefault("cross", "github_user", "")
	crossRepoPath := config.GetOrDefault("cross", "cross_repo_path", "")

	if githubUser == "" || crossRepoPath == "" {
		utils.Warn("Cross-build not configured. Please run 'juice --setup cross' first.")
		os.Exit(1)
	}

	utils.Step("Initiating Cloud Cross-Build for " + platform)
	utils.EnsureGHCLILoggedIn()

	if !utils.DirExists(crossRepoPath) {
		utils.Info("Cross-build repository not found at " + crossRepoPath + ". Attempting to restore...")
		utils.RunCmd("gh", "repo", "clone", githubUser+"/juice-cross", crossRepoPath)
	}

	utils.Info("Syncing codebase to cloud builder...")
	files, _ := os.ReadDir(crossRepoPath)
	for _, f := range files {
		if f.Name() == ".git" || f.Name() == ".github" {
			continue
		}
		_ = os.RemoveAll(filepath.Join(crossRepoPath, f.Name()))
	}

	if err := utils.CopyDirectory(".", crossRepoPath); err != nil {
		utils.Fatal("Failed to sync code to cross-repo", err)
	}

	_ = os.RemoveAll(filepath.Join(crossRepoPath, ".native"))
	_ = os.RemoveAll(filepath.Join(crossRepoPath, "build"))
	_ = os.RemoveAll(filepath.Join(crossRepoPath, "temps"))

	if err := sanitizeCrossRepoGoMod(crossRepoPath); err != nil {
		utils.Warn("Failed to sanitize go.mod for cross build: " + err.Error())
	}

	utils.BuildFrontend()

	origWd, _ := os.Getwd()
	_ = os.Chdir(crossRepoPath)
	utils.RunCmd("git", "add", ".")
	_ = exec.Command("git", "commit", "-m", "chore: automated cross-build sync").Run()
	utils.RunCmd("git", "push", "origin", "main")
	_ = os.Chdir(origWd)

	utils.WaitForActionFinish(crossRepoPath)

	utils.Println()
	utils.Info("Downloading built bindings from GitHub Release...")
	iosNativePath := filepath.Join(".native", "ios")
	_ = os.MkdirAll(iosNativePath, 0755)

	zipPath := filepath.Join(iosNativePath, "Sweetjuice.xcframework.zip")
	releaseURL := fmt.Sprintf("https://github.com/%s/juice-cross/releases/latest/download/Sweetjuice.xcframework.zip", githubUser)

	if err := utils.DownloadFile(releaseURL, zipPath); err != nil {
		utils.Fatal("Failed to download built framework from release", err)
	}

	utils.Info("Extracting framework...")
	if err := utils.UnzipTarget(zipPath, iosNativePath); err != nil {
		utils.Fatal("Failed to extract framework", err)
	}
	_ = os.Remove(zipPath)

	utils.Success("Cloud build integration successful. Launching on local hardware...")
	ios.RunPipeline()
}

func ExecuteSetup(target string) {
	if target != "cross" {
		utils.Error("Unknown setup target '" + target + "'. Available: 'cross'")
		os.Exit(1)
	}

	utils.Step("Setting up Cross-Compilation Environment")
	utils.EnsureGHCLILoggedIn()

	out, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		utils.Fatal("Failed to get GitHub user info", err)
	}
	githubUser := strings.TrimSpace(string(out))
	utils.Printf("Detected GitHub user: %s\n", githubUser)

	home, _ := os.UserHomeDir()
	crossRepoPath := filepath.Join(home, ".sweetjuice", "juice-cross")

	utils.Println("")
	utils.Info("Step 1: Forking sweet-juice/juice-cross...")
	utils.Println("Please ensure you have manually forked https://github.com/sweet-juice/juice-cross to your account.")

	utils.Println("")
	utils.Info("Step 2: Cloning your fork locally...")
	if !utils.DirExists(filepath.Dir(crossRepoPath)) {
		_ = os.MkdirAll(filepath.Dir(crossRepoPath), 0755)
	}

	if utils.DirExists(crossRepoPath) {
		utils.Printf("Repository already exists at %s. Resetting to latest remote state...\n", crossRepoPath)
		origWd, _ := os.Getwd()
		_ = os.Chdir(crossRepoPath)
		utils.RunCmd("git", "fetch", "origin")
		utils.RunCmd("git", "reset", "--hard", "origin/main")
		_ = os.Chdir(origWd)
	} else {
		utils.RunCmd("gh", "repo", "clone", githubUser+"/juice-cross", crossRepoPath)
	}

	utils.Println("")
	utils.Info("Step 3: Updating project configuration...")
	if utils.FileExists("config.ini") {
		data, _ := os.ReadFile("config.ini")
		content := string(data)
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "github_user =") {
				lines[i] = "github_user = " + githubUser
				break
			}
		}
		for i, line := range lines {
			if strings.HasPrefix(line, "cross_repo_path =") {
				lines[i] = "cross_repo_path = " + crossRepoPath
				break
			}
		}
		content = strings.Join(lines, "\n")
		_ = os.WriteFile("config.ini", []byte(content), 0644)
		utils.Success("Successfully updated config.ini")
	} else {
		utils.Warn("config.ini not found in current directory. Configuration skipped.")
	}

	utils.Println("")
	utils.Success("Cross-compilation setup complete!")
}

func sanitizeCrossRepoGoMod(repoPath string) error {
	goModPath := filepath.Join(repoPath, "go.mod")
	if !utils.FileExists(goModPath) {
		return fmt.Errorf("go.mod not found in cross repo")
	}

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "replace github.com/sweet-juice/sweetjuice =>") {
			result = append(result, "\treplace github.com/sweet-juice/sweetjuice => .")
			continue
		}
		result = append(result, line)
	}

	sanitized := strings.Join(result, "\n")
	return os.WriteFile(goModPath, []byte(sanitized), 0644)
}

func ManagePlugin(action, pluginRepo string) {
	if !utils.DirExists(".plugins") || !utils.FileExists("go.mod") {
		utils.Error("You must execute plugin commands from the root of a juice project directory.")
		os.Exit(1)
	}

	if action == "add" {
		utils.Step("Installing Plugin: " + pluginRepo)
		utils.RunCmd("go", "get", pluginRepo)

		out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", pluginRepo).Output()
		if err != nil {
			utils.Error("Could not resolve source location for module: " + pluginRepo)
			os.Exit(1)
		}
		goModSrc := strings.TrimSpace(string(out))

		androidSrc := filepath.Join(goModSrc, "android")
		if utils.DirExists(androidSrc) {
			utils.Info("Syncing Android native directory...")
			if err := utils.CopyDirectory(androidSrc, filepath.Join(".plugins", "android")); err != nil {
				utils.Fatal("Failed to sync Android native directory", err)
			}
		}

		iosSrc := filepath.Join(goModSrc, "ios")
		if utils.DirExists(iosSrc) {
			utils.Info("Syncing iOS native directory...")
			if err := utils.CopyDirectory(iosSrc, filepath.Join(".plugins", "ios")); err != nil {
				utils.Fatal("Failed to sync iOS native directory", err)
			}
		}
		utils.Success("Plugin " + pluginRepo + " added successfully!")
	} else if action == "remove" {
		utils.Step("Removing Plugin: " + pluginRepo)
		pluginDirname := filepath.Base(pluginRepo)
		_ = os.RemoveAll(filepath.Join(".plugins", "android", pluginDirname))
		_ = os.RemoveAll(filepath.Join(".plugins", "ios", pluginDirname))

		utils.RunCmd("go", "mod", "edit", "-droprequire="+pluginRepo)
		utils.RunCmd("go", "mod", "tidy")
	}
}
