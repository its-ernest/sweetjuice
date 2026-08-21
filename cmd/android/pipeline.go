package android

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

var (
	GomobileTarget = "android/arm64"
	AndroidAPI     = "21"
	AarName        = "sweetjuice.aar"
	CleanOutput    = "true"
)

type AndroidManifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	PackageName string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

type Application struct {
	Activities []Activity `xml:"activity"`
}

type Activity struct {
	Name          string         `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters []IntentFilter `xml:"intent-filter"`
}

type IntentFilter struct {
	Actions    []Action   `xml:"action"`
	Categories []Category `xml:"category"`
}

type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

func SetupAndroidLocalProperties(targetDir string) {
	sdkPath := GetAndroidSDKPath()
	if sdkPath != "" {
		androidDir := filepath.Join(targetDir, ".native", "android")
		propsPath := filepath.Join(androidDir, "local.properties")
		if !utils.DirExists(androidDir) {
			_ = os.MkdirAll(androidDir, 0755)
		}
		content := fmt.Sprintf("sdk.dir=%s\n", filepath.ToSlash(sdkPath))
		_ = os.WriteFile(propsPath, []byte(content), 0644)
		utils.Info("[android] SDK found at " + sdkPath + ". Updated local.properties")
	} else {
		utils.Warn("[android] Could not locate Android SDK automatically. Please set ANDROID_HOME.")
	}
}

func GetAndroidSDKPath() string {
	if path := os.Getenv("ANDROID_HOME"); path != "" && utils.DirExists(path) {
		return path
	}
	if path := os.Getenv("ANDROID_SDK_ROOT"); path != "" && utils.DirExists(path) {
		return path
	}

	home, _ := os.UserHomeDir()
	var paths []string

	switch runtime.GOOS {
	case "darwin":
		paths = []string{filepath.Join(home, "Library", "Android", "sdk")}
	case "linux":
		paths = []string{
			filepath.Join(home, "Android", "Sdk"),
			filepath.Join(home, "android-sdk"),
			"/usr/lib/android-sdk",
			"/opt/android-sdk",
		}
	case "windows":
		paths = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk"),
		}
	}

	for _, p := range paths {
		if utils.DirExists(p) {
			return p
		}
	}

	return ""
}

func findApksigner() string {
	if utils.CommandExists("apksigner") {
		return "apksigner"
	}

	sdkPath := GetAndroidSDKPath()
	if sdkPath == "" {
		return ""
	}

	buildToolsDir := filepath.Join(sdkPath, "build-tools")
	entries, err := os.ReadDir(buildToolsDir)
	if err != nil {
		return ""
	}

	type versionInfo struct {
		name string
		mtime int64
	}
	var versions []versionInfo
	for _, e := range entries {
		if e.IsDir() {
			info, _ := e.Info()
			versions = append(versions, versionInfo{e.Name(), info.ModTime().UnixNano()})
		}
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i].mtime > versions[j].mtime })

	for _, v := range versions {
		candidate := filepath.Join(buildToolsDir, v.name, "apksigner")
		if runtime.GOOS == "windows" {
			candidate += ".bat"
		}
		if utils.FileExists(candidate) {
			return candidate
		}
	}

	return ""
}

func signApkTo(srcPath, destPath, keystorePath, alias, keystorePass, keyPassword string) error {
	apksigner := findApksigner()
	if apksigner == "" {
		utils.Warn("apksigner not found, skipping APK signing")
		return fmt.Errorf("apksigner not found")
	}

	if !utils.FileExists(keystorePath) {
		utils.Warn("Keystore not found at " + keystorePath + ", skipping APK signing")
		return fmt.Errorf("keystore not found")
	}

	utils.Info("Signing " + filepath.Base(srcPath) + " -> " + filepath.Base(destPath))
	cmd := exec.Command(apksigner, "sign",
		"--ks", keystorePath,
		"--ks-key-alias", alias,
		"--ks-pass", "pass:"+keystorePass,
		"--key-pass", "pass:"+keyPassword,
		"--out", destPath,
		srcPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		utils.Warn("APK signing failed: " + err.Error())
		return err
	}

	utils.Info("Signed APK written to " + destPath)
	return nil
}

func ValidateAndroidEnvironment() {
	sdkPath := GetAndroidSDKPath()
	if sdkPath == "" {
		utils.Fatal("Android SDK not found", fmt.Errorf("could not locate SDK in environment or default locations. Please install Android Studio or set ANDROID_HOME"))
	}

	if utils.DirExists(filepath.Join(".native", "android")) {
		SetupAndroidLocalProperties(".")
	}

	utils.EnsureGoMobileTools()
}

func applyConfig() {
	config := utils.LoadConfig()
	target := config.GetOrDefault("build", "android_target", "arm64")
	if !strings.Contains(target, "/") {
		GomobileTarget = "android/" + target
	} else {
		GomobileTarget = target
	}
	AndroidAPI = config.GetOrDefault("android", "min_api", "21")
	AarName = config.GetOrDefault("android", "aar_name", "sweetjuice.aar")
}

func applyAppConfig() {
	config := utils.LoadConfig()
	appName := config.GetOrDefault("app", "name", "Sweet Juice")
	packageName := config.GetOrDefault("app", "package", "com.sweetjuice.app")
	versionCode := config.GetOrDefault("build", "versionCode", "1")
	versionName := config.GetOrDefault("build", "version", "1.0")
	minApi := config.GetOrDefault("android", "min_api", "23")

	settingsPath := filepath.Join(".native", "android", "settings.gradle")
	if utils.FileExists(settingsPath) {
		data, _ := os.ReadFile(settingsPath)
		updated := strings.Replace(string(data), "rootProject.name = \"SweetJuice Mobile\"", "rootProject.name = \""+appName+"\"", 1)
		_ = os.WriteFile(settingsPath, []byte(updated), 0644)
	}

	buildGradlePath := filepath.Join(".native", "android", "app", "build.gradle")
	if utils.FileExists(buildGradlePath) {
		data, _ := os.ReadFile(buildGradlePath)
		updated := string(data)
		updated = regexp.MustCompile(`namespace\s+'[^']+'`).ReplaceAllString(updated, "namespace '"+packageName+"'")
		updated = regexp.MustCompile(`applicationId\s+"[^"]+"`).ReplaceAllString(updated, "applicationId \""+packageName+"\"")
		updated = regexp.MustCompile(`minSdk\s+\d+`).ReplaceAllString(updated, "minSdk "+minApi)
		updated = regexp.MustCompile(`versionCode\s+\d+`).ReplaceAllString(updated, "versionCode "+versionCode)
		updated = regexp.MustCompile(`versionName\s+"[^"]+"`).ReplaceAllString(updated, "versionName \""+versionName+"\"")
		_ = os.WriteFile(buildGradlePath, []byte(updated), 0644)
	}

	manifestPath := filepath.Join(".native", "android", "app", "src", "main", "AndroidManifest.xml")
	if utils.FileExists(manifestPath) {
		data, _ := os.ReadFile(manifestPath)
		updated := strings.Replace(string(data), "package=\"com.sweetjuice.app\"", "package=\""+packageName+"\"", 1)
		_ = os.WriteFile(manifestPath, []byte(updated), 0644)
	}

	stringsPath := filepath.Join(".native", "android", "app", "src", "main", "res", "values", "strings.xml")
	if utils.FileExists(stringsPath) {
		data, _ := os.ReadFile(stringsPath)
		updated := strings.Replace(string(data), "<string name=\"app_name\">Sweet Juice</string>", "<string name=\"app_name\">"+appName+"</string>", 1)
		_ = os.WriteFile(stringsPath, []byte(updated), 0644)
	}

	javaBaseDir := filepath.Join(".native", "android", "app", "src", "main", "java", "com", "sweetjuice", "app")
	if utils.DirExists(javaBaseDir) {
		files, _ := os.ReadDir(javaBaseDir)
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".java") {
				filePath := filepath.Join(javaBaseDir, f.Name())
				data, _ := os.ReadFile(filePath)
				updated := strings.Replace(string(data), "package com.sweetjuice.app;", "package "+packageName+";", 1)
				_ = os.WriteFile(filePath, []byte(updated), 0644)
			}
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

	ext := filepath.Ext(iconPath)
	iconName := "ic_launcher" + ext
	roundIconName := "ic_launcher_round" + ext
	foregroundName := "ic_launcher_foreground" + ext

	mipmapDirs, _ := filepath.Glob(filepath.Join(".native", "android", "app", "src", "main", "res", "mipmap-*"))
	for _, dir := range mipmapDirs {
		base := filepath.Base(dir)
		if base == "mipmap-anydpi-v26" {
			continue
		}

		entries, _ := os.ReadDir(dir)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if (strings.HasPrefix(name, "ic_launcher.") || strings.HasPrefix(name, "ic_launcher_round.") || strings.HasPrefix(name, "ic_launcher_foreground.")) && (strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".webp")) {
				_ = os.Remove(filepath.Join(dir, name))
			}
		}

		_ = os.WriteFile(filepath.Join(dir, iconName), data, 0644)
		_ = os.WriteFile(filepath.Join(dir, roundIconName), data, 0644)
		_ = os.WriteFile(filepath.Join(dir, foregroundName), data, 0644)
	}

	adaptiveIconDir := filepath.Join(".native", "android", "app", "src", "main", "res", "mipmap-anydpi-v26")
	if utils.DirExists(adaptiveIconDir) {
		launcherXML := filepath.Join(adaptiveIconDir, "ic_launcher.xml")
		if utils.FileExists(launcherXML) {
			data, _ := os.ReadFile(launcherXML)
			updated := strings.ReplaceAll(string(data), "@drawable/ic_launcher_foreground", "@mipmap/ic_launcher_foreground")
			_ = os.WriteFile(launcherXML, []byte(updated), 0644)
		}
		roundXML := filepath.Join(adaptiveIconDir, "ic_launcher_round.xml")
		if utils.FileExists(roundXML) {
			data, _ := os.ReadFile(roundXML)
			updated := strings.ReplaceAll(string(data), "@drawable/ic_launcher_foreground", "@mipmap/ic_launcher_foreground")
			_ = os.WriteFile(roundXML, []byte(updated), 0644)
		}
	}

	manifestPath := filepath.Join(".native", "android", "app", "src", "main", "AndroidManifest.xml")
	if utils.FileExists(manifestPath) {
		data, _ := os.ReadFile(manifestPath)
		updated := strings.Replace(string(data), "android:icon=\"@mipmap/ic_juice\"", "android:icon=\"@mipmap/ic_launcher\"", 1)
		updated = strings.Replace(updated, "android:roundIcon=\"@mipmap/ic_juice_round\"", "android:roundIcon=\"@mipmap/ic_launcher_round\"", 1)
		_ = os.WriteFile(manifestPath, []byte(updated), 0644)
	}

	backgroundDrawable := filepath.Join(".native", "android", "app", "src", "main", "res", "drawable", "ic_launcher_background.xml")
	if utils.FileExists(backgroundDrawable) {
		_ = os.WriteFile(backgroundDrawable, []byte(`<?xml version="1.0" encoding="utf-8"?>
<vector xmlns:android="http://schemas.android.com/apk/res/android"
    android:width="108dp"
    android:height="108dp"
    android:viewportWidth="108"
    android:viewportHeight="108">
    <path
        android:fillColor="#FFFFFF"
        android:pathData="M0,0h108v108h-108z" />
</vector>`), 0644)
	}

	utils.Info("App icon replaced with " + iconPath)
}

func RefreshPipeline() {
	ValidateAndroidEnvironment()
	applyConfig()
	applyAppConfig()

	assetsDir := filepath.Join(".native", "android", "app", "src", "main", "assets")
	if err := utils.CopyAppAssets(assetsDir); err != nil {
		utils.Warn("Failed to copy app_assets: " + err.Error())
	}

	utils.BuildFrontend()

	outputPath := filepath.Join(".native", "android", "app", "libs")
	targetJavaSrcDir := filepath.Join(".native", "android", "app", "src", "main", "java")
	stagingPluginsDir := filepath.Join(".plugins", "android")

	_ = os.MkdirAll(outputPath, 0755)

	if CleanOutput == "true" {
		files, _ := filepath.Glob(filepath.Join(outputPath, "*.aar"))
		for _, f := range files {
			_ = os.Remove(f)
		}
	}

	utils.Info("Building " + AarName + " for target " + GomobileTarget + "...")
	utils.RunCmd("gomobile", "bind", "-target="+GomobileTarget, "-androidapi="+AndroidAPI, "-o", filepath.Join(outputPath, AarName), ".")

	if utils.DirExists(stagingPluginsDir) && !utils.DirEmpty(stagingPluginsDir) {
		_ = os.MkdirAll(targetJavaSrcDir, 0755)
		if err := utils.CopyDirectory(stagingPluginsDir, targetJavaSrcDir); err != nil {
			utils.Fatal("Failed syncing plugin package trees inside Android workspace", err)
		}
	}
}

func BuildPipeline(mode string) {
	ValidateAndroidEnvironment()
	applyConfig()
	androidDir := filepath.Join(".native", "android")
	if !utils.DirExists(androidDir) {
		utils.Error("Native Android path layout missing.")
		os.Exit(1)
	}

	gradleCmd := "./gradlew"
	if utils.IsWindowsHost() {
		gradleCmd = "gradlew.bat"
	}

	origWd, _ := os.Getwd()
	targetDir := filepath.Join(origWd, ".build", "android", mode)

	config := utils.LoadConfig()
	keystore := config.GetOrDefault("android", "keystore", "")
	keyAlias := config.GetOrDefault("android", "alias", config.GetOrDefault("android", "key_alias", ""))
	keystorePass := config.GetOrDefault("android", "keystore_pass", config.GetOrDefault("android", "key_password", ""))
	keyPassword := config.GetOrDefault("android", "key_password", "")

	_ = os.Chdir(androidDir)
	defer func() { _ = os.Chdir(origWd) }()

	var targetTask string
	switch mode {
	case "debug":
		targetTask = "assembleDebug"
	case "release":
		targetTask = "assembleRelease"
	case "bundle":
		targetTask = "bundleRelease"
	}

	utils.RunCmd(gradleCmd, targetTask)

	buildOutputBase := "app/build/outputs"
	_ = os.MkdirAll(targetDir, 0755)

	// Clean previous artifacts in target to avoid stale/duplicate files
	entries, _ := os.ReadDir(targetDir)
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".apk") || strings.HasSuffix(entry.Name(), ".aab")) {
			_ = os.Remove(filepath.Join(targetDir, entry.Name()))
		}
	}

	var apkPaths []string

		var copied int
		switch mode {
		case "debug":
			apkDir := filepath.Join(buildOutputBase, "apk", "debug")
			if utils.DirExists(apkDir) {
				files, _ := filepath.Glob(filepath.Join(apkDir, "*.apk"))
				for _, f := range files {
					base := filepath.Base(f)
					_ = utils.CopyFile(f, filepath.Join(targetDir, base))
					apkPaths = append(apkPaths, filepath.Join(targetDir, base))
					copied++
				}
			}
		case "release":
			apkDir := filepath.Join(buildOutputBase, "apk", "release")
			if utils.DirExists(apkDir) {
				files, _ := filepath.Glob(filepath.Join(apkDir, "*.apk"))
				for _, f := range files {
					base := filepath.Base(f)
					_ = utils.CopyFile(f, filepath.Join(targetDir, base))
					apkPaths = append(apkPaths, filepath.Join(targetDir, base))
					copied++
				}
			}
			aabDir := filepath.Join(buildOutputBase, "bundle", "release")
			if utils.DirExists(aabDir) {
				files, _ := filepath.Glob(filepath.Join(aabDir, "*.aab"))
				for _, f := range files {
					_ = utils.CopyFile(f, filepath.Join(targetDir, filepath.Base(f)))
					copied++
				}
			}
		case "bundle":
			aabDir := filepath.Join(buildOutputBase, "bundle", "release")
			if utils.DirExists(aabDir) {
				files, _ := filepath.Glob(filepath.Join(aabDir, "*.aab"))
				for _, f := range files {
					_ = utils.CopyFile(f, filepath.Join(targetDir, filepath.Base(f)))
					copied++
				}
			}
		}

	if copied > 0 {
		utils.Info(fmt.Sprintf("Copied %d build artifact(s) to %s", copied, targetDir))
	} else {
		utils.Warn("No build artifacts found to copy to " + targetDir)
	}

	if mode != "debug" && keystore != "" && keyAlias != "" && keystorePass != "" && keyPassword != "" {
		resolvedKeystore := keystore
		if !filepath.IsAbs(resolvedKeystore) {
			resolvedKeystore = filepath.Join(origWd, resolvedKeystore)
		}
		utils.Info("Signing " + fmt.Sprintf("%d", len(apkPaths)) + " APK(s) with keystore " + resolvedKeystore)
		for _, copiedPath := range apkPaths {
			if !strings.HasSuffix(copiedPath, ".apk") {
				continue
			}
			base := filepath.Base(copiedPath)
			ext := filepath.Ext(base)
			nameWithoutExt := strings.TrimSuffix(base, ext)
			signedName := strings.TrimSuffix(nameWithoutExt, "-unsigned") + ext
			signedPath := filepath.Join(targetDir, signedName)
			if err := signApkTo(copiedPath, signedPath, resolvedKeystore, keyAlias, keystorePass, keyPassword); err != nil {
				utils.Error("Signing failed for " + base + ": " + err.Error())
			}
		}
	} else if mode == "debug" {
		utils.Info("Skipping APK signing for debug build.")
	} else {
		utils.Warn("Skipping APK signing: keystore config missing (need [android] keystore, alias, keystore_pass, key_password in config.ini)")
	}
}

func RunPipeline() {
	ValidateAndroidEnvironment()

	adbTool := "adb"
	if _, err := exec.LookPath("adb"); err != nil {
		sdkPath := GetAndroidSDKPath()
		adbPath := filepath.Join(sdkPath, "platform-tools", "adb")
		if runtime.GOOS == "windows" {
			adbPath += ".exe"
		}
		if _, err := os.Stat(adbPath); err == nil {
			adbTool = adbPath
		} else {
			utils.Error("'adb' tool missing and not found in SDK platform-tools.")
			os.Exit(1)
		}
	}

	apkPath := filepath.Join(".native", "android", "app", "build", "outputs", "apk", "debug", "app-debug.apk")
	if !utils.FileExists(apkPath) {
		utils.Error("APK not found at " + apkPath + ". Did you build it?")
		os.Exit(1)
	}

	utils.Info("Installing APK to device...")
	utils.RunCmd(adbTool, "install", "-r", apkPath)

	manifestPath := filepath.Join(".native", "android", "app", "src", "main", "AndroidManifest.xml")
	packageName, launcherActivity, err := parseManifestDetails(manifestPath)
	if err != nil {
		return
	}

	if packageName == "" {
		config := utils.LoadConfig()
		packageName = config.GetOrDefault("app", "package", "com.sweetjuice.app")
	}

	if strings.HasPrefix(launcherActivity, ".") {
		launcherActivity = packageName + launcherActivity
	}

	utils.Info("Launching application " + packageName + "/" + launcherActivity + "...")
	utils.RunCmd(adbTool, "shell", "am", "start", "-n", packageName+"/"+launcherActivity)
}

func parseManifestDetails(manifestPath string) (string, string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", "", err
	}
	var manifest AndroidManifest
	if err := xml.Unmarshal(data, &manifest); err != nil {
		return "", "", err
	}

	packageName := manifest.PackageName
	launcherActivity := ""
	activityLoop:
	for _, act := range manifest.Application.Activities {
		for _, filter := range act.IntentFilters {
			for _, action := range filter.Actions {
				if action.Name == "android.intent.action.MAIN" {
					launcherActivity = act.Name
					break activityLoop
				}
			}
		}
	}
	return packageName, launcherActivity, nil
}
