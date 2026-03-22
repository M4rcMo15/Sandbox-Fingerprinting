"""File-based XSS injection."""

import os
import tempfile
from datetime import datetime
from .base import BaseInjector
from ..payloads import XSSPayload


class FileInjector(BaseInjector):
    """Injects XSS payloads into various file types."""
    
    def inject(self, payload: XSSPayload):
        """Create multiple files with XSS payloads."""
        temp_dir = tempfile.gettempdir()
        timestamp = datetime.now().isoformat()
        
        # HTML files
        self._create_html_file(temp_dir, payload, timestamp)
        
        # Text files
        self._create_text_file(temp_dir, payload, timestamp)
        
        # JSON files
        self._create_json_file(temp_dir, payload, timestamp)
        
        # XML files
        self._create_xml_file(temp_dir, payload, timestamp)
        
        # Markdown files
        self._create_markdown_file(temp_dir, payload, timestamp)
    
    def _create_html_file(self, temp_dir: str, payload: XSSPayload, timestamp: str):
        """Create HTML file with XSS payload."""
        html_content = f"""<!DOCTYPE html>
<html>
<head>
    <title>XSS Test {payload.id}</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>XSS Payload Test</h1>
    <p>Payload ID: {payload.id}</p>
    <p>Type: {payload.type}</p>
    <p>Timestamp: {timestamp}</p>
    <div id="xss-test">{payload.content}</div>
    <script>
        var xssPayload = '{payload.content}';
        console.log('XSS Test:', xssPayload);
    </script>
    <!-- XSS Comment: {payload.content} -->
    <div data-xss="{payload.content}" style="display:none">{payload.content}</div>
</body>
</html>"""
        
        file_path = os.path.join(temp_dir, f"xss_test_{payload.id}.html")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(html_content)
    
    def _create_text_file(self, temp_dir: str, payload: XSSPayload, timestamp: str):
        """Create text file with XSS payload."""
        txt_content = f"""XSS Test Log
Payload ID: {payload.id}
Type: {payload.type}
Timestamp: {timestamp}
Payload: {payload.content}
Test: {payload.content}
Description: {payload.description}
"""
        
        file_path = os.path.join(temp_dir, f"xss_log_{payload.id}.txt")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(txt_content)
    
    def _create_json_file(self, temp_dir: str, payload: XSSPayload, timestamp: str):
        """Create JSON file with XSS payload."""
        # Note: Not using json.dumps to preserve raw XSS payload
        json_content = f"""{{
  "test_id": "{payload.id}",
  "type": "{payload.type}",
  "payload": "{payload.content}",
  "xss_test": "{payload.content}",
  "timestamp": "{timestamp}",
  "description": "{payload.description}"
}}"""
        
        file_path = os.path.join(temp_dir, f"xss_config_{payload.id}.json")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(json_content)
    
    def _create_xml_file(self, temp_dir: str, payload: XSSPayload, timestamp: str):
        """Create XML file with XSS payload."""
        xml_content = f"""<?xml version="1.0" encoding="UTF-8"?>
<XSSTest>
  <ID>{payload.id}</ID>
  <Type>{payload.type}</Type>
  <Payload>{payload.content}</Payload>
  <Test>{payload.content}</Test>
  <Timestamp>{timestamp}</Timestamp>
  <Description>{payload.description}</Description>
</XSSTest>"""
        
        file_path = os.path.join(temp_dir, f"xss_report_{payload.id}.xml")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(xml_content)
    
    def _create_markdown_file(self, temp_dir: str, payload: XSSPayload, timestamp: str):
        """Create Markdown file with XSS payload."""
        md_content = f"""# XSS Test {payload.id}

## Payload
{payload.content}

## Test
{payload.content}

## Details
- **Type**: {payload.type}
- **Timestamp**: {timestamp}
- **Description**: {payload.description}

## Raw Payload
```
{payload.content}
```
"""
        
        file_path = os.path.join(temp_dir, f"xss_readme_{payload.id}.md")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(md_content)
