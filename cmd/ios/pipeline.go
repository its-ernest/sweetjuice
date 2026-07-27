package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

var (
	CleanOutput      = "true"
	IOSTarget        = "ios"
	XFrameworkName   = "Sweetjuice.xcframework"
)

func applyConfig() {
	config := utils.LoadConfig()
	target := config.GetOrDefault("build", "ios_target", "ios")
	if target != "ios" && !strings.Contains(target, "/") {
		IOSTarget = "ios/" + target
	} else {
		IOSTarget = target
	}
	XFrameworkName = "Sweetjuice.xcframework"
}

func ValidateIOSEnvironment() {
	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool missing", fmt.Errorf("please install xtool: https://github.com/mizage/xtool"))
	}

	utils.EnsureGoMobileTools()
}

func ScaffoldProject(name string) {
	utils.Info("Scaffolding iOS project '" + name + "' using xtool...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		_ = os.MkdirAll(iosDir, 0755)
	}

	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool is required for iOS project scaffolding", fmt.Errorf("xtool command not found"))
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	utils.Info("Running xtool scaffolding sequence...")
	utils.RunCmd("xtool", "init", "--name", name)

	utils.Success("iOS project scaffolded successfully.")
}

func RefreshPipeline() {
	ValidateIOSEnvironment()
	applyConfig()

	utils.BuildFrontend()

	utils.Info("Starting iOS toolchain refresh and Go compilation...")

	iosDir := filepath.Join("native", "ios")
	targetSrcDir := filepath.Join(iosDir, "Sources", "Plugins")
	stagingPluginsDir := filepath.Join(".plugins", "ios")
	outputPath := filepath.Join(iosDir, XFrameworkName)

	if !utils.DirExists(iosDir) {
		utils.Error("Native iOS path layout missing. Ensure you have the iOS template files in native/ios.")
		os.Exit(1)
	}

	if CleanOutput == "true" {
		if utils.DirExists(outputPath) {
			utils.Info("Cleaning previous Go bindings...")
			_ = os.RemoveAll(outputPath)
		}
	}

	utils.Info("Generating Go bindings (XCFramework) for iOS via gomobile...")
	if runtime.GOOS != "darwin" {
		utils.Warn("[ios] Note: Local iOS binding on Linux/Windows requires a specialized gomobile toolchain.")
		utils.Debug("It is recommended to use 'juice --run-cross ios' for GitHub Action based builds.")
	}

	utils.RunCmd("gomobile", "bind", "-target="+IOSTarget, "-o", outputPath, ".")

	if utils.DirExists(stagingPluginsDir) && !utils.DirEmpty(stagingPluginsDir) {
		utils.Info("Syncing iOS native plugins...")
		_ = os.MkdirAll(targetSrcDir, 0755)
		if err := utils.CopyDirectory(stagingPluginsDir, targetSrcDir); err != nil {
			utils.Fatal("Failed syncing plugin package trees inside iOS workspace", err)
		}
	}

	utils.Info("iOS Platform Refresh complete. Bindings generated in " + outputPath)
}

func BuildPipeline(mode string) {
	ValidateIOSEnvironment()
	applyConfig()
	utils.Info("Building iOS application in mode: " + mode + " via xtool...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		utils.Error("Native iOS path layout missing.")
		os.Exit(1)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	config := "debug"
	if mode == "release" || mode == "bundle" {
		config = "release"
	}

	utils.Info("Executing xtool build sequence (" + config + ")...")
	utils.RunCmd("xtool", "dev", "build", "--configuration", config)

	utils.Success("iOS application built successfully.")
}

func RunPipeline() {
	ValidateIOSEnvironment()
	utils.Info("Preparing device deployment via xtool...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		utils.Error("Native iOS path layout missing.")
		os.Exit(1)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	utils.Info("Deploying to connected iOS device...")
	utils.RunCmd("xtool", "dev", "run")
}
