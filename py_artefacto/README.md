# XSS Sandbox Audit

A Python tool for auditing sandbox security by injecting XSS payloads into various system vectors.

## Description

This tool tests the security of malware analysis sandboxes by injecting XSS payloads into multiple system components including:

- File names and content
- Registry keys and values
- Process arguments and environment variables
- Network traffic (HTTP headers)
- Windows Event Logs
- Scheduled tasks
- Service descriptions
- Clipboard
- Debug output
- And many more vectors

## Installation

```bash
pip install xss-sandbox-audit
```

## Usage

### Command Line

```bash
# Run full XSS audit
xss-audit

# Run specific injection vectors
xss-audit --vectors files,registry,network

# Use custom XSS payload
xss-audit --payload '<script>alert("XSS")</script>'
```

### Python API

```python
from xss_audit import XSSAuditor

# Create auditor instance
auditor = XSSAuditor()

# Run all injection vectors
auditor.inject_all()

# Run specific vectors
auditor.inject_files()
auditor.inject_registry()
auditor.inject_network()
```

## Features

- **Multiple Injection Vectors**: Tests 20+ different system components
- **Customizable Payloads**: Use default or custom XSS payloads
- **Windows-Specific**: Optimized for Windows sandbox environments
- **Stealth Mode**: Runs with hidden windows to avoid detection
- **Comprehensive Coverage**: Targets all major sandbox analysis points

## Injection Vectors

1. **Files**: HTML, TXT, JSON, XML, MD files with XSS payloads
2. **Registry**: Multiple registry keys with XSS values
3. **Processes**: Process arguments containing XSS
4. **PowerShell**: Various PowerShell commands with XSS
5. **CMD**: Command prompt commands with XSS
6. **Event Logs**: Windows Event Log entries with XSS
7. **Environment Variables**: System environment variables with XSS
8. **Network**: HTTP headers with XSS payloads
9. **Shortcuts**: LNK files with XSS in descriptions
10. **Scheduled Tasks**: Windows tasks with XSS
11. **Services**: Service descriptions with XSS
12. **Clipboard**: Clipboard content with XSS
13. **Debug Output**: OutputDebugString with XSS
14. **WMI**: WMI queries with XSS context
15. **DNS**: DNS queries with XSS markers

## Requirements

- Python 3.8+
- Windows OS
- Administrator privileges (recommended for full functionality)

## Warning

This tool is designed for security research and testing purposes only. Use only in controlled environments with proper authorization.

## License

MIT License
