"""Basic usage examples for XSS Sandbox Audit."""

from xss_audit import XSSAuditor, XSSPayload


def example_full_audit():
    """Run a full XSS audit with all vectors."""
    print("Example 1: Full XSS Audit")
    print("-" * 50)
    
    auditor = XSSAuditor()
    auditor.inject_all(delay=0.2)
    
    print("\nFull audit completed!\n")


def example_specific_vectors():
    """Run audit on specific vectors only."""
    print("Example 2: Specific Vectors")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Inject only into files and registry
    print("[+] Injecting into files...")
    auditor.inject_files()
    
    print("[+] Injecting into registry...")
    auditor.inject_registry()
    
    print("\nSpecific vector audit completed!\n")


def example_custom_payload():
    """Use a custom XSS payload."""
    print("Example 3: Custom Payload")
    print("-" * 50)
    
    # Create custom payload
    custom_payload = XSSPayload(
        id="custom_001",
        type="custom-alert",
        content='<script>alert("Custom XSS")</script>',
        vector="all",
        description="Custom alert-based XSS payload"
    )
    
    # Create auditor with custom payload
    auditor = XSSAuditor(payloads=[custom_payload])
    
    # Run audit
    auditor.inject_all(delay=0.1)
    
    print("\nCustom payload audit completed!\n")


def example_multiple_custom_payloads():
    """Use multiple custom XSS payloads."""
    print("Example 4: Multiple Custom Payloads")
    print("-" * 50)
    
    payloads = [
        XSSPayload(
            id="custom_001",
            type="script-tag",
            content='<script>console.log("Test 1")</script>',
            vector="files",
            description="Script tag payload"
        ),
        XSSPayload(
            id="custom_002",
            type="img-onerror",
            content='<img src=x onerror=alert("Test 2")>',
            vector="files",
            description="Image onerror payload"
        ),
        XSSPayload(
            id="custom_003",
            type="svg-onload",
            content='<svg onload=alert("Test 3")>',
            vector="all",
            description="SVG onload payload"
        ),
    ]
    
    auditor = XSSAuditor(payloads=payloads)
    auditor.inject_all(delay=0.1)
    
    print("\nMultiple custom payloads audit completed!\n")


def example_targeted_injection():
    """Target specific injection methods."""
    print("Example 5: Targeted Injection")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Only inject into network and environment
    print("[+] Injecting into network traffic...")
    auditor.inject_network()
    
    print("[+] Injecting into environment variables...")
    auditor.inject_environment()
    
    print("[+] Injecting into PowerShell...")
    auditor.inject_powershell()
    
    print("\nTargeted injection completed!\n")


if __name__ == "__main__":
    print("=" * 50)
    print("XSS Sandbox Audit - Usage Examples")
    print("=" * 50)
    print()
    
    # Run examples
    example_full_audit()
    example_specific_vectors()
    example_custom_payload()
    example_multiple_custom_payloads()
    example_targeted_injection()
    
    print("=" * 50)
    print("All examples completed!")
    print("=" * 50)
