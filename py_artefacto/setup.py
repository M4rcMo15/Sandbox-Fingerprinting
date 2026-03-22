from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="xss-sandbox-audit",
    version="1.0.0",
    author="Security Research Team",
    author_email="security@example.com",
    description="XSS payload injection tool for sandbox security auditing",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourusername/xss-sandbox-audit",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Information Technology",
        "Topic :: Security",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Operating System :: Microsoft :: Windows",
    ],
    python_requires=">=3.8",
    install_requires=[
        "requests>=2.28.0",
    ],
    entry_points={
        "console_scripts": [
            "xss-audit=xss_audit.cli:main",
        ],
    },
)
