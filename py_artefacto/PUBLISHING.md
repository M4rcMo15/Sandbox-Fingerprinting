# Publishing to PyPI

This guide explains how to publish the `xss-sandbox-audit` package to PyPI.

## Prerequisites

1. Install required tools:
```bash
pip install build twine
```

2. Create accounts:
   - PyPI: https://pypi.org/account/register/
   - TestPyPI (optional, for testing): https://test.pypi.org/account/register/

## Build the Package

1. Clean previous builds:
```bash
rm -rf build/ dist/ *.egg-info
```

2. Build the package:
```bash
python -m build
```

This creates:
- `dist/xss-sandbox-audit-1.0.0.tar.gz` (source distribution)
- `dist/xss_sandbox_audit-1.0.0-py3-none-any.whl` (wheel distribution)

## Test on TestPyPI (Optional but Recommended)

1. Upload to TestPyPI:
```bash
python -m twine upload --repository testpypi dist/*
```

2. Test installation:
```bash
pip install --index-url https://test.pypi.org/simple/ xss-sandbox-audit
```

3. Test the package:
```bash
xss-audit --list-vectors
```

## Publish to PyPI

1. Upload to PyPI:
```bash
python -m twine upload dist/*
```

2. Enter your PyPI credentials when prompted.

3. Verify the upload at: https://pypi.org/project/xss-sandbox-audit/

## Install from PyPI

Once published, users can install with:
```bash
pip install xss-sandbox-audit
```

## Version Updates

To publish a new version:

1. Update version in:
   - `setup.py`
   - `pyproject.toml`
   - `xss_audit/__init__.py`

2. Clean and rebuild:
```bash
rm -rf build/ dist/ *.egg-info
python -m build
```

3. Upload new version:
```bash
python -m twine upload dist/*
```

## Using API Tokens (Recommended)

For better security, use API tokens instead of passwords:

1. Generate token at: https://pypi.org/manage/account/token/

2. Create `~/.pypirc`:
```ini
[pypi]
username = __token__
password = pypi-AgEIcHlwaS5vcmc...your-token-here...
```

3. Upload:
```bash
python -m twine upload dist/*
```

## Troubleshooting

### Package name already exists
If the package name is taken, change it in:
- `setup.py` (name parameter)
- `pyproject.toml` ([project] name)

### Upload fails
- Check credentials
- Verify package builds correctly: `python -m build`
- Check for errors: `python -m twine check dist/*`

### Import errors after installation
- Ensure package structure is correct
- Check `__init__.py` files exist in all directories
- Verify `packages` in `setup.py` includes all modules

## Best Practices

1. Always test on TestPyPI first
2. Use semantic versioning (MAJOR.MINOR.PATCH)
3. Update CHANGELOG.md for each release
4. Tag releases in git: `git tag v1.0.0`
5. Use API tokens instead of passwords
6. Keep README.md up to date
7. Include comprehensive examples

## Resources

- PyPI: https://pypi.org/
- TestPyPI: https://test.pypi.org/
- Python Packaging Guide: https://packaging.python.org/
- Twine Documentation: https://twine.readthedocs.io/
