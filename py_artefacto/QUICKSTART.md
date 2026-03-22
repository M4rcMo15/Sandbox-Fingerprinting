# Quick Start Guide

Get started with XSS Sandbox Audit in 5 minutes!

## Installation

```bash
pip install xss-sandbox-audit
```

## Basic Usage

### Command Line

Run a full XSS audit:
```bash
xss-audit
```

List available vectors:
```bash
xss-audit --list-vectors
```

Test specific vectors:
```bash
xss-audit --vectors files,registry,network
```

Use custom payload:
```bash
xss-audit --payload '<script>alert("XSS")</script>'
```

### Python API

```python
from xss_audit import XSSAuditor

# Create auditor and run full audit
auditor = XSSAuditor()
auditor.inject_all()
```

## Common Scenarios

### Test Specific Sandbox

```python
from xss_audit import XSSAuditor

auditor = XSSAuditor()

# Test vectors that Any.Run analyzes
auditor.inject_processes()
auditor.inject_network()
auditor.inject_registry()
auditor.inject_files()
```

### Custom Payload

```python
from xss_audit import XSSAuditor, XSSPayload

payload = XSSPayload(
    id="custom_001",
    type="custom",
    content='<script>alert("Test")</script>',
    vector="all",
    description="My custom payload"
)

auditor = XSSAuditor(payloads=[payload])
auditor.inject_all()
```

### Stealth Mode

```python
from xss_audit import XSSAuditor
import time

auditor = XSSAuditor()

# Inject slowly
auditor.inject_environment()
time.sleep(5)
auditor.inject_debug()
time.sleep(5)
auditor.inject_clipboard()
```

## Available Vectors

| Vector | Description |
|--------|-------------|
| files | HTML, TXT, JSON, XML, MD files |
| registry | Windows Registry keys |
| processes | Process arguments |
| powershell | PowerShell commands |
| cmd | CMD commands |
| eventlogs | Windows Event Logs |
| environment | Environment variables |
| network | HTTP headers |
| shortcuts | LNK files |
| tasks | Scheduled tasks |
| services | Service descriptions |
| clipboard | Clipboard content |
| debug | OutputDebugString |
| wmi | WMI queries |
| dns | DNS queries |

## Next Steps

- Read the full [README.md](README.md)
- Check out [examples/](examples/)
- Review [PUBLISHING.md](PUBLISHING.md) to publish your own version
- See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute

## Need Help?

- Check the [examples/](examples/) directory
- Read the full documentation in [README.md](README.md)
- Open an issue on GitHub

## Warning

⚠️ This tool is for authorized security testing only. Use responsibly!
