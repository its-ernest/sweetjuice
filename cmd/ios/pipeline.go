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

func sanitizeSwiftName(name string) string {
	words := strings.Fields(name)
	result := ""
	for _, w := range words {
		if len(w) > 0 {
			result += strings.ToUpper(w[:1]) + w[1:]
		}
	}
	if result == "" {
		result = "App"
	}
	return result
}

func applyAppConfig() {
	config := utils.LoadConfig()
	appName := config.GetOrDefault("app", "name", "Sweet Juice")
	packageName := config.GetOrDefault("app", "package", "com.sweetjuice.app")
	sanitizedName := sanitizeSwiftName(appName)

	xtoolPath := filepath.Join(".native", "ios", "xtool.yml")
	if utils.FileExists(xtoolPath) {
		data, _ := os.ReadFile(xtoolPath)
		updated := strings.Replace(string(data), "bundleID: com.sweetjuice.app", "bundleID: "+packageName, 1)
		_ = os.WriteFile(xtoolPath, []byte(updated), 0644)
	}

	packageSwiftPath := filepath.Join(".native", "ios", "Package.swift")
	if utils.FileExists(packageSwiftPath) {
		data, _ := os.ReadFile(packageSwiftPath)
		updated := strings.ReplaceAll(string(data), "\"GenericApp\"", "\""+sanitizedName+"\"")
		_ = os.WriteFile(packageSwiftPath, []byte(updated), 0644)
	}

	oldSourceDir := filepath.Join(".native", "ios", "Sources", "GenericApp")
	newSourceDir := filepath.Join(".native", "ios", "Sources", sanitizedName)
	if utils.DirExists(oldSourceDir) && sanitizedName != "GenericApp" {
		_ = os.Rename(oldSourceDir, newSourceDir)
	}

	if utils.DirExists(newSourceDir) {
		swiftAppPath := filepath.Join(newSourceDir, sanitizedName+"App.swift")
		if utils.FileExists(swiftAppPath) {
			data, _ := os.ReadFile(swiftAppPath)
			updated := strings.Replace(string(data), "struct GenericAppApp", "struct "+sanitizedName+"App", 1)
			updated = strings.Replace(updated, "GenericAppApp: App", sanitizedName+"App: App", 1)
			_ = os.WriteFile(swiftAppPath, []byte(updated), 0644)
		}
	}

	applyIconConfig()
}

func applyIconConfig() {
	config := utils.LoadConfig()
	iconPath := config.GetOrDefault("app", "icon", "")
	if iconPath == "" {
		return
	}

	if !utils.FileExists(iconPath) {
		utils.Warn("App icon not found at " + iconPath + ", skipping icon replacement")
		return
	}

	data, err := os.ReadFile(iconPath)
	if err != nil {
		utils.Warn("Failed to read app icon: " + err.Error())
		return
	}

	iosResourcesDir := filepath.Join(".native", "ios", "Resources")
	if !utils.DirExists(iosResourcesDir) {
		_ = os.MkdirAll(iosResourcesDir, 0755)
	}
	iconDest := filepath.Join(iosResourcesDir, "AppIcon.png")
	_ = os.WriteFile(iconDest, data, 0644)

	plistPaths := []string{
		filepath.Join(".native", "ios", "xtool", "App.app", "Info.plist"),
	}
	for _, plistPath := range plistPaths {
		if !utils.FileExists(plistPath) {
			continue
		}
		plistData, _ := os.ReadFile(plistPath)
		updated := string(plistData)
		if !strings.Contains(updated, "CFBundleIconFile") {
			updated = strings.Replace(updated, "<key>CFBundleInfoDictionaryVersion</key>", "<key>CFBundleIconFile</key>\n\t<string>AppIcon</string>\n\t<key>CFBundleInfoDictionaryVersion</key>", 1)
			_ = os.WriteFile(plistPath, []byte(updated), 0644)
		}
	}

	utils.Info("iOS app icon copied to " + iconDest)
}

func ValidateIOSEnvironment() {
	if !utils.CommandExists("xtool") {
		utils.Fatal("xtool missing", fmt.Errorf("please install xtool: https://github.com/mizage/xtool"))
	}

	utils.EnsureGoMobileTools()
}

func ScaffoldProject(name string) {
	utils.Info("Scaffolding iOS project '" + name + "' using xtool...")

	iosDir := filepath.Join(".native", "ios")
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
	applyAppConfig()

	iosDir := filepath.Join(".native", "ios")
	assetsDir := filepath.Join(iosDir, "Resources")
	if err := utils.CopyAppAssets(assetsDir); err != nil {
		utils.Warn("Failed to copy app_assets: " + err.Error())
	}

	utils.BuildFrontend()

	targetSrcDir := filepath.Join(iosDir, "Sources", "Plugins")
	stagingPluginsDir := filepath.Join(".plugins", "ios")
	outputPath := filepath.Join(iosDir, XFrameworkName)

	if !utils.DirExists(iosDir) {
		utils.Error("Native iOS path layout missing. Ensure you have the iOS template files in .native/ios.")
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

	iosDir := filepath.Join(".native", "ios")
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

	iosDir := filepath.Join(".native", "ios")
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
