# Contributing to XSS Sandbox Audit

Thank you for your interest in contributing to XSS Sandbox Audit! This document provides guidelines for contributing to the project.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- Clear description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- System information (OS, Python version)
- Any relevant logs or screenshots

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:
- Clear description of the enhancement
- Use cases and benefits
- Possible implementation approach
- Any relevant examples

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`pytest`)
6. Update documentation as needed
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/xss-sandbox-audit.git
cd xss-sandbox-audit
```

2. Create a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install development dependencies:
```bash
pip install -e .
pip install -r requirements-dev.txt
```

4. Run tests:
```bash
pytest
```

## Coding Standards

### Python Style Guide

- Follow PEP 8
- Use type hints where appropriate
- Write docstrings for all public functions and classes
- Keep functions focused and small
- Use meaningful variable names

### Code Formatting

Use `black` for code formatting:
```bash
black xss_audit/
```

### Linting

Use `flake8` for linting:
```bash
flake8 xss_audit/
```

### Type Checking

Use `mypy` for type checking:
```bash
mypy xss_audit/
```

## Testing

### Writing Tests

- Write tests for all new functionality
- Use descriptive test names
- Follow the Arrange-Act-Assert pattern
- Mock external dependencies

### Running Tests

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=xss_audit

# Run specific test file
pytest tests/test_payloads.py

# Run specific test
pytest tests/test_payloads.py::TestPayloads::test_get_default_payloads
```

## Documentation

### Docstring Format

Use Google-style docstrings:

```python
def function_name(param1: str, param2: int) -> bool:
    """
    Brief description of function.
    
    Longer description if needed.
    
    Args:
        param1: Description of param1
        param2: Description of param2
        
    Returns:
        Description of return value
        
    Raises:
        ValueError: Description of when this is raised
    """
    pass
```

### Updating Documentation

- Update README.md for user-facing changes
- Update CHANGELOG.md for all changes
- Add examples for new features
- Update docstrings for modified functions

## Adding New Injection Vectors

To add a new injection vector:

1. Create a new injector in `xss_audit/injectors/`:
```python
from .base import BaseInjector
from ..payloads import XSSPayload

class NewInjector(BaseInjector):
    """Injects XSS payloads into new vector."""
    
    def inject(self, payload: XSSPayload):
        """Inject payload into new vector."""
        # Implementation here
        pass
```

2. Add to `xss_audit/injectors/__init__.py`:
```python
from .new_injector import NewInjector

__all__ = [
    # ... existing injectors ...
    "NewInjector",
]
```

3. Add to `XSSAuditor` in `xss_audit/auditor.py`:
```python
def __init__(self, ...):
    # ... existing injectors ...
    self.new_injector = NewInjector()

def inject_new_vector(self):
    """Inject XSS payloads into new vector."""
    for payload in self.payloads:
        if payload.vector in ["all", "new_vector"]:
            self.new_injector.inject(payload)
```

4. Add to CLI in `xss_audit/cli.py`

5. Write tests in `tests/test_new_injector.py`

6. Update documentation

## Commit Messages

Use clear and descriptive commit messages:

```
Add DNS injection vector

- Implement DNSInjector class
- Add nslookup and ipconfig commands
- Add tests for DNS injection
- Update documentation
```

Format:
- First line: Brief summary (50 chars or less)
- Blank line
- Detailed description with bullet points

## Release Process

1. Update version in:
   - `setup.py`
   - `pyproject.toml`
   - `xss_audit/__init__.py`

2. Update CHANGELOG.md

3. Create git tag:
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

4. Build and publish:
```bash
python -m build
python -m twine upload dist/*
```

## Questions?

If you have questions, feel free to:
- Open an issue
- Contact the maintainers
- Check existing documentation

Thank you for contributing!
