# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-22

### Added
- Initial release of XSS Sandbox Audit
- 15 different injection vectors:
  - File-based injection (HTML, TXT, JSON, XML, MD)
  - Windows Registry injection
  - Process argument injection
  - PowerShell command injection
  - CMD command injection
  - Windows Event Log injection
  - Environment variable injection
  - Network traffic injection (HTTP headers)
  - Windows Shortcut (LNK) injection
  - Scheduled Task injection
  - Windows Service injection
  - Clipboard injection
  - Debug output injection (OutputDebugString)
  - WMI query injection
  - DNS query injection
- Command-line interface with `xss-audit` command
- Python API for programmatic usage
- Default XSS payloads with base64 encoding
- Contextual payload variations
- Custom payload support
- Comprehensive documentation and examples
- MIT License

### Features
- Stealth mode with hidden windows
- Configurable delay between injections
- Vector selection (all or specific)
- Multiple payload support
- Error handling and silent failures
- Cross-vector payload injection

## [Unreleased]

### Planned
- Additional injection vectors
- Callback server support for payload tracking
- Enhanced reporting capabilities
- Configuration file support
- Logging and verbose mode
- Payload templates
- Integration with CI/CD pipelines
