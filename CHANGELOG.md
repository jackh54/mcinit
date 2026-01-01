# Changelog

All notable changes to mcinit will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of mcinit
- Cross-platform support (Windows, macOS, Linux)
- Server initialization (`mcinit init`)
- Server lifecycle management (`start`, `stop`, `restart`, `status`)
- Log viewing and filtering (`mcinit logs`)
- Support for multiple server types:
  - Vanilla (Mojang)
  - Paper
  - Purpur
  - Folia
  - Velocity
  - Waterfall
- Automatic Java detection and validation
- JVM flags presets (Aikar's, minimal, custom)
- Startup script generation (Unix shell, PowerShell, CMD)
- Reproducible configuration via `mcinit.json`
- .gitignore integration
- Cache management for server jars
- Homebrew, Scoop, and Winget installation support

## [0.1.0] - 2025-01-XX

### Added
- Initial MVP release

