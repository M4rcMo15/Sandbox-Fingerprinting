"""Command-line interface for XSS Sandbox Audit."""

import argparse
import sys
from .auditor import XSSAuditor
from .payloads import XSSPayload, get_default_payloads


def main():
    """Main CLI entry point."""
    parser = argparse.ArgumentParser(
        description="XSS Sandbox Audit - Security testing tool for sandbox environments"
    )
    
    parser.add_argument(
        "--vectors",
        type=str,
        help="Comma-separated list of vectors to test (e.g., files,registry,network). Default: all",
        default="all"
    )
    
    parser.add_argument(
        "--payload",
        type=str,
        help="Custom XSS payload to use instead of defaults",
        default=None
    )
    
    parser.add_argument(
        "--delay",
        type=float,
        help="Delay in seconds between injection batches (default: 0.1)",
        default=0.1
    )
    
    parser.add_argument(
        "--list-vectors",
        action="store_true",
        help="List all available injection vectors"
    )
    
    args = parser.parse_args()
    
    if args.list_vectors:
        print_available_vectors()
        return 0
    
    # Create auditor
    if args.payload:
        custom_payload = XSSPayload(
            id="custom_001",
            type="custom",
            content=args.payload,
            vector="all",
            description="Custom user-provided payload"
        )
        auditor = XSSAuditor(payloads=[custom_payload])
    else:
        auditor = XSSAuditor()
    
    # Run audit
    if args.vectors == "all":
        auditor.inject_all(delay=args.delay)
    else:
        run_specific_vectors(auditor, args.vectors, args.delay)
    
    return 0


def print_available_vectors():
    """Print all available injection vectors."""
    vectors = [
        ("files", "File names and content (HTML, TXT, JSON, XML, MD)"),
        ("registry", "Windows Registry keys and values"),
        ("processes", "Process arguments"),
        ("powershell", "PowerShell commands"),
        ("cmd", "CMD commands"),
        ("eventlogs", "Windows Event Logs"),
        ("environment", "Environment variables"),
        ("network", "Network traffic (HTTP headers)"),
        ("shortcuts", "Windows shortcuts (LNK files)"),
        ("tasks", "Scheduled tasks"),
        ("services", "Service descriptions"),
        ("clipboard", "Clipboard content"),
        ("debug", "Debug output (OutputDebugString)"),
        ("wmi", "WMI queries"),
        ("dns", "DNS queries"),
    ]
    
    print("\nAvailable Injection Vectors:")
    print("-" * 60)
    for name, description in vectors:
        print(f"  {name:15} - {description}")
    print()


def run_specific_vectors(auditor: XSSAuditor, vectors_str: str, delay: float):
    """Run specific injection vectors."""
    vector_map = {
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
    
    vectors = [v.strip() for v in vectors_str.split(",")]
    
    print("[*] Starting XSS audit with selected vectors...")
    
    for vector in vectors:
        if vector not in vector_map:
            print(f"[-] Unknown vector: {vector}")
            continue
        
        print(f"[+] Injecting into {vector}...")
        try:
            vector_map[vector]()
            import time
            time.sleep(delay)
        except Exception as e:
            print(f"[-] Error injecting into {vector}: {e}")
    
    print("[*] XSS audit complete!")


if __name__ == "__main__":
    sys.exit(main())
