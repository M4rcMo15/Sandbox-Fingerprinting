# XSS Sandbox Audit - Examples

This directory contains example scripts demonstrating various usage patterns of the XSS Sandbox Audit tool.

## Basic Usage Examples

### `basic_usage.py`

Demonstrates fundamental usage patterns:

1. **Full Audit**: Run all injection vectors
2. **Specific Vectors**: Target specific injection points
3. **Custom Payload**: Use a single custom XSS payload
4. **Multiple Custom Payloads**: Use multiple custom payloads
5. **Targeted Injection**: Inject into selected vectors

Run with:
```bash
python examples/basic_usage.py
```

## Advanced Usage Examples

### `advanced_usage.py`

Demonstrates advanced techniques:

1. **Staggered Injection**: Inject with delays for stealth
2. **Payload Variations**: Test multiple XSS payload types
3. **Targeted Sandbox Test**: Focus on sandbox analysis points
4. **Persistence Vectors**: Target persistence mechanisms
5. **Evasion-Focused**: Use low-visibility vectors
6. **Comprehensive Report**: Generate detailed test reports

Run with:
```bash
python examples/advanced_usage.py
```

## Example Scenarios

### Testing a Specific Sandbox

```python
from xss_audit import XSSAuditor

# Create auditor
auditor = XSSAuditor()

# Test vectors that Any.Run analyzes
auditor.inject_processes()  # Process tree
auditor.inject_network()    # Network traffic
auditor.inject_registry()   # Registry changes
auditor.inject_files()      # File operations
```

### Custom Payload Testing

```python
from xss_audit import XSSAuditor, XSSPayload

# Create custom payload
payload = XSSPayload(
    id="custom_001",
    type="custom",
    content='<script>alert("XSS")</script>',
    vector="all",
    description="Custom test payload"
)

# Test with custom payload
auditor = XSSAuditor(payloads=[payload])
auditor.inject_all()
```

### Stealth Testing

```python
from xss_audit import XSSAuditor
import time

auditor = XSSAuditor()

# Inject slowly to avoid detection
vectors = [
    auditor.inject_environment,
    auditor.inject_debug,
    auditor.inject_clipboard,
]

for inject_func in vectors:
    inject_func()
    time.sleep(5)  # 5 second delay
```

### Comprehensive Testing

```python
from xss_audit import XSSAuditor

# Test everything with default payloads
auditor = XSSAuditor()
auditor.inject_all(delay=0.2)
```

## Tips

1. **Start Small**: Begin with specific vectors before running full audit
2. **Use Delays**: Add delays between injections for stealth
3. **Custom Payloads**: Test with payloads specific to your target
4. **Monitor Results**: Check sandbox reports to see which vectors work
5. **Iterate**: Refine your approach based on results

## Common Use Cases

### Security Research
Test how different sandboxes handle XSS in various contexts.

### Sandbox Evaluation
Evaluate sandbox XSS filtering and sanitization capabilities.

### Penetration Testing
Test XSS vulnerabilities in analysis environments.

### Red Team Operations
Assess detection and response to XSS injection attempts.

## Safety Notes

- Only use in authorized testing environments
- Do not use against production systems
- Respect legal and ethical boundaries
- Follow responsible disclosure practices

## More Information

See the main README.md for:
- Installation instructions
- Full API documentation
- Command-line usage
- Contributing guidelines
