# CLI Reference

Compact reference and useful examples for the `juice` CLI.

Synopsis

```
Sweet Juice Toolchain CLI (juice)
Usage:
  juice --new <project_name>        Create a fresh project from the template
  juice --refresh <platform>        Run platform sync: 'android' or 'ios'
  juice --build <platform> <mode>   Sync environment and compile app binary (debug/release/bundle)
  juice --run <platform>            Compile, install, and execute application via ADB
  juice --add <plugin-url>          Install a native Go/Mobile plugin
  juice --remove <plugin-url>       Uninstall a native Go/Mobile plugin
```

Basic examples

```bash
# Create a new project
juice --new my-app

# Build debug APK
juice --build android debug

# Build and run (the CLI performs sync automatically)
juice --build android debug && juice --run android

# Build release bundle (AAB)
juice --build android bundle

# Add external plugin
juice --add github.com/author/plugin

# Remove plugin
juice --remove github.com/author/plugin
```

When to use `--refresh`

- The CLI automatically runs sync during `build`/`run`. Use `juice --refresh android` only when you want to force a rebuild of the generated AAR/JAR artifacts (for debugging or template changes).
- `--refresh` now reads `config.ini` and applies the configured `app.name`, `app.package`, `build.version`, `build.versionCode`, and `android.min_api` to the Android template, and `app.name`, `app.package`, and `build.version` to the iOS template.
- `--refresh` also copies `app_assets/` into the platform asset directories (`native/android/app/src/main/assets/` and `native/ios/Resources/`).
- If `config.ini [app] icon` is set, `--refresh` replaces the Android launcher icons in all `mipmap-*` directories and updates `AndroidManifest.xml` accordingly, and copies the icon to `native/ios/Resources/AppIcon.png`.

Environment notes

- Ensure Android SDK and NDK are installed and available to the environment (Android Studio recommended).
- If builds fail due to Gradle, open `native/android/` in Android Studio and inspect the Gradle wrapper and SDK/NDK settings.

CI / scripts

- For CI pipelines prefer explicit commands: `juice --refresh android && juice --build android release` to keep steps visible.
- After a successful Android build, the CLI copies the produced APK/AAB artifacts into `.build/android/debug/` or `.build/android/release/` depending on the build mode.
- If `config.ini [android] keystore`, `alias`, `keystore_pass`, and `key_password` are set and the keystore file exists, the CLI will sign the copied APK using `apksigner` after copying, keeping both the signed (`app-release.apk`) and unsigned (`app-release-unsigned.apk`) versions.
